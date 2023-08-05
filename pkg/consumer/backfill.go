package consumer

import (
	"context"
	"net/http"
	"sync"
	"time"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/repo"
	"github.com/ericvolp12/bsky-experiments/pkg/consumer/store/store_queries"
	"github.com/ipfs/go-cid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type BackfillRepoStatus struct {
	RepoDid     string
	Seq         int64
	State       string
	EventBuffer []*comatproto.SyncSubscribeRepos_Commit
	lk          sync.Mutex
}

type RecordJob struct {
	RecordPath string
	NodeCid    cid.Cid
}

type RecordResult struct {
	RecordPath string
	Error      error
}

func (c *Consumer) BackfillProcessor(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "BackfillProcessor")
	defer span.End()

	log := c.Logger.With("source", "backfill_main")
	log.Info("starting backfill processor")

	// Create a semaphore with a capacity of 50
	sem := make(chan struct{}, 50)

	for {
		select {
		case <-ctx.Done():
			log.Info("stopping backfill processor")
			return
		default:
		}

		// Get the next backfill
		var backfill *BackfillRepoStatus
		c.statusLock.RLock()
		for _, b := range c.BackfillStatus {
			b.lk.Lock()
			if b.State == "enqueued" {
				backfill = b
				b.State = "in_progress"
				b.lk.Unlock()
				break
			}
			b.lk.Unlock()
		}
		c.statusLock.RUnlock()

		if backfill == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		sem <- struct{}{} // Block until there is a slot in the semaphore
		go func(b *BackfillRepoStatus) {
			// Process the backfill
			c.ProcessBackfill(ctx, b.RepoDid)
			backfillJobsProcessed.WithLabelValues(c.SocketURL).Inc()
			<-sem // Release a slot in the semaphore when the goroutine finishes
		}(backfill)
	}
}

func (c *Consumer) ProcessBackfill(ctx context.Context, repoDID string) {
	ctx, span := tracer.Start(ctx, "ProcessBackfill")
	defer span.End()

	log := c.Logger.With("source", "backfill", "repo", repoDID)
	log.Infof("processing backfill for %s", repoDID)

	var url = "https://bsky.social/xrpc/com.atproto.sync.getCheckout?did=" + repoDID

	// GET and CAR decode the body
	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   120 * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Errorf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Accept", "application/vnd.ipld.car")

	c.SyncLimiter.Wait(ctx)

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error sending request: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Error response: %v", resp.StatusCode)
		return
	}

	r, err := repo.ReadRepoFromCar(ctx, resp.Body)
	if err != nil {
		log.Errorf("Error reading repo: %v", err)
		return
	}

	numRecords := 0
	numRoutines := 50
	recordJobs := make(chan RecordJob, numRoutines)
	recordResults := make(chan RecordResult, numRoutines)

	wg := sync.WaitGroup{}

	// Producer routine
	go func() {
		defer close(recordJobs)
		r.ForEach(ctx, "", func(recordPath string, nodeCid cid.Cid) error {
			numRecords++
			recordJobs <- RecordJob{RecordPath: recordPath, NodeCid: nodeCid}
			return nil
		})
	}()

	// Consumer routines
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range recordJobs {
				recordCid, rec, err := r.GetRecord(ctx, job.RecordPath)
				if err != nil {
					log.Errorf("Error getting record: %v", err)
					recordResults <- RecordResult{RecordPath: job.RecordPath, Error: err}
					continue
				}

				// Verify that the record cid matches the cid in the event
				if recordCid != job.NodeCid {
					log.Errorf("mismatch in record and op cid: %s != %s", recordCid, job.NodeCid)
					recordResults <- RecordResult{RecordPath: job.RecordPath, Error: err}
					continue
				}

				_, err = c.HandleCreateRecord(ctx, repoDID, job.RecordPath, rec)
				if err != nil {
					log.Errorf("failed to handle create record: %+v", err)
				}

				backfillRecordsProcessed.WithLabelValues(c.SocketURL).Inc()
				recordResults <- RecordResult{RecordPath: job.RecordPath, Error: err}
			}
		}()
	}

	resultWG := sync.WaitGroup{}
	resultWG.Add(1)
	// Handle results
	go func() {
		defer resultWG.Done()
		for result := range recordResults {
			if result.Error != nil {
				log.Errorf("Error processing record %s: %v", result.RecordPath, result.Error)
			}
		}
	}()

	wg.Wait()
	close(recordResults)
	resultWG.Wait()

	log.Infof("processed %d records", numRecords)

	c.statusLock.RLock()
	bf := c.BackfillStatus[repoDID]
	c.statusLock.RUnlock()

	// Update the backfill record
	bf.lk.Lock()
	err = c.Store.Queries.UpdateRepoBackfillRecord(ctx, store_queries.UpdateRepoBackfillRecordParams{
		Repo:         repoDID,
		LastBackfill: time.Now(),
		SeqStarted:   bf.Seq,
		State:        "complete",
	})
	if err != nil {
		log.Errorf("failed to update repo backfill record: %+v", err)
	}

	// Update the backfill status
	bf.State = "complete"
	bf.lk.Unlock()

	bufferedEventsProcessed := 0
	log.Infof("processing %d buffered events", len(bf.EventBuffer))
	// Playback the buffered events
	for _, evt := range bf.EventBuffer {
		err = c.HandleRepoCommit(ctx, evt)
		if err != nil {
			log.Errorf("failed to handle repo commit: %+v", err)
		}
		backfillEventsBuffered.WithLabelValues(c.SocketURL).Dec()
		bufferedEventsProcessed++
	}

	log.Infof("processed %d buffered events", bufferedEventsProcessed)
}
