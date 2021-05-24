# topshot-alert
Monitor NBA TopShots purchase/listing events using Flow, Go, Angular, Google Cloud Functions, MySQL.

## Background

[NBA TopShots](https://nbatopshot.com/) is a NFT collectible website which sells NBA game 'moments' in packs.  It has a marketplace which allows collectors to buy/sell the moments they own.  The NFT transactions are recorded on the [Flow blockchain](https://www.onflow.org/) so every moment sale and listing are recorded on the Flow blockchain.

This utility uses the Flow SDK to monitor purchase/listing events in realtime to attempt to better understand the prices of moments or to detect a specific moment the user may be interested in buying.

## purchase-events 

This program queries the Flow blockchain for events related to the NBA TopShots contract.  This includes events when a moment is listed or purchased.  The events can be displayed to the console or be configured to upload them to a MySQL database.  The program is written in Go.  See [purchase-events](purchase-events/README.md)

```
PS ..\purchase-events> go run .\query_loop.go   
{"N":"Jarrett Allen","P":1,"SN":22766,"T":"Base Set","level":"info","msg":"Purchased","time":"2021-05-21T19:24:32-04:00"}
{"N":"Jrue Holiday","P":2,"SN":9978,"T":"Base Set","level":"info","msg":"Purchased","time":"2021-05-21T19:24:32-04:00"}
{"N":"Karl-Anthony Towns","P":3,"SN":20242,"T":"Base Set","level":"info","msg":"Purchased","time":"2021-05-21T19:24:32-04:00"}
{"N":"Karl-Anthony Towns","P":21,"SN":17612,"T":"Hustle and Show","level":"info","msg":"Listed","time":"2021-05-21T19:24:33-04:00"}
```

## topshot-functions/go-db-functions
 
The is a Google Cloud Function which takes the events created in purchase-events and creates a simple query API which other programs can use to see purchase/listing events without running a purchase-events collector.  The program is written in Go.  See [topshot-functions/go-db-functions](topshot-functions/go-db-functions/README.md)

## topshot-alert-ui

The is a simple Angular web UI application which uses the API provided by topshot-functions to display the latest purchase/events in a web browser.  It has functions like filtering by specific moments/sets.  See [topshot-alert-ui](topshot-alert-ui/README.md)

## Purpose

I did this largely as a learning expirement.  I had no previous expirence in Go programming or blockchains.


## Acknowledgments

This was inspired and contains elements from https://medium.com/@eric.ren_51534/polling-nba-top-shot-p2p-market-purchase-events-from-flow-blockchain-using-flow-go-sdk-3ec80119e75f


## Contact

If you have any questions about the code or other exciting opportunities please contact nochopfoundation@gmail.com  I'm a 20+ year developer working mostly with Java.  I've worked on essentially the same product/team for many years so it was fun to try and learn a few technologies.  