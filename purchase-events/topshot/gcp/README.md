# Deploying on Google Cloud

## Compute Instance
I provisioned a small Ubuntu VM installing
- go
- mysql command-line utility to administor/test the GCP SQL instance running MySQL

## Long Running Event Collector
The main purpose is to run the long running process which 
- Queries the Flow blockchain for NBA TopShot contract events
- Uploads the database with the events

I have a simple bash script which runs an infinte loop.  Ideally, you could create a systemd service and provide cloud alerting in-case of errors.

```
#!/bin/bash
export DB_CONNECTION="root:*******@tcp(10.122.49.6:3306)/topshotdev"
export COLLECTOR_ID=1
while true
do
        echo "START RUN LOOP" >> console
        echo $(date) >> console
        go run ./query_loop.go 1000 >> console 2>> console.err 
        sleep 60
done
```

launching via
```
./query_run &
```