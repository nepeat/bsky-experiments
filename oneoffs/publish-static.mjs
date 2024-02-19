import dotenv from 'dotenv'
import atprotoAPI from '@atproto/api';
import fs from 'fs/promises'

const FEED_GEN_TYPE = "app.bsky.feed.generator";
const feedGenDid = "did:web:feed-api.owo.me";

async function publishFeed(agent, feed, recordData) {
    let avatarRef;
    if (recordData.avatar) {
        let encoding;
        if (recordData.avatar.endsWith('png')) {
            encoding = 'image/png'
        } else if (recordData.avatar.endsWith('jpg') || recordData.avatar.endsWith('jpeg')) {
            encoding = 'image/jpeg'
        } else {
            throw new Error('expected png or jpeg')
        }
        const img = await fs.readFile(recordData.avatar)
        const blobRes = await agent.api.com.atproto.repo.uploadBlob(img, {
            encoding,
        })
        avatarRef = blobRes.data.blob
    }
    
    let record = {
        repo: agent.session?.did ?? '',
        collection: FEED_GEN_TYPE,
        rkey: feed,
        record: {
            did: feedGenDid,
            displayName: recordData.displayName,
            description: recordData.description,
            avatar: avatarRef,
            createdAt: new Date().toISOString(),
        },
    };
    console.log(record);
    await agent.api.com.atproto.repo.putRecord(record)
}

async function deleteFeed(agent, feed) {
    console.log(`Deleting ${feed}...`)
    await agent.api.com.atproto.repo.deleteRecord({
        repo: agent.session?.did ?? '',
        collection: FEED_GEN_TYPE,
        rkey: feed,
    })
}

const run = async () => {
    dotenv.config()
    
    let auth = process.env.ATP_AUTH;
    let handle = auth.split(':')[0];
    let password = auth.split(':')[1];
    
    const feeds = {
        // "dot_com": {
        //     displayName: "TLD Feed - .com",
        //     description: "Firehose content from the .com TLD",
        // },
        // "dot_net": {
        //     displayName: "TLD Feed - .net",
        //     description: "Firehose content from the .net TLD",
        // },
        // "dot_org": {
        //     displayName: "TLD Feed - .org",
        //     description: "Firehose content from the .org TLD",
        // },
        // "dot_me": {
        //     displayName: "TLD Feed - .me",
        //     description: "Firehose content from the .me TLD",
        // },
        // "dot_io": {
        //     displayName: "TLD Feed - .io",
        //     description: "Firehose content from the .io TLD",
        // },
        // "avoid5": {
        //     displayName: "AVoid5 @ BSky",
        //     description: "r/AVoid5 but on BSky",
        // },
        // "novowels": {
        //     displayName: "No Vow*ls",
        //     description: "why????",
        // },
        "testfeed": {
            displayName: "test feed, please ignore",
            description: "test",
        },
    }
    
    // (Optional) The path to an image to be used as your feed's avatar
    // Ex: ~/path/to/avatar.jpeg
    const avatar = ''
    
    // only update this if in a test environment
    const agent = new atprotoAPI.BskyAgent({ service: 'https://bsky.social' })
    await agent.login({ identifier: handle, password })
    
    for (const [feed, recordData] of Object.entries(feeds)) {
        console.log(`Publishing ${feed}...`)
        
        if (process.env.DELETE) {
            await deleteFeed(agent, feed);
        } else {
            await publishFeed(agent, feed, recordData)
        }
    };
    console.log('All done 🎉')
}

run()