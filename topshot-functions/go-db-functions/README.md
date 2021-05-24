# Simple Query API Using Go and Google Cloud Functions
This provides a simple backend API to the NBA TopShot purchase/listing events collected by the ***purchase-events*** application.

## GET /plays
Returns the static metadata for each moment.

## GET /status/collectors/recent
Returns recent status on each ***purchase-events*** collector.  This can be used to determine if the collectors are currently popluating purchase/listing events.

## GET /stream/momentEvents/from/HEAD
Returns the most recent purchase/listing events.  This should be the first call to start streaming events.  After establishing the latest block height, calls to ***/stream/momentEvents/from/HEAD/to/blockHeight/{BlockHeight}*** should be made to avoid receiving duplicate events.

## GET /stream/momentEvents/from/HEAD/to/blockHeight/{BlockHeight}
Returns the most recent purchase/listing events after BlockHeight.  This can be called every 20-30 seconds to create a stream of events.

## GET /momentEvents/play/:PlayId/set/:SetId/from/HEAD
Returns fairly recent events specific to a play moment in a set.

# Deploying to Google Cloud Functions
It is fairly easy to deploy to Google Cloud Functions.  I created a small Google SQL instance and configured the functions to access the database.

```
gcloud functions deploy GetPlayDataHTTP --runtime go113 --trigger-http --allow-unauthenticated --set-env-vars DB_USER=root --set-env-vars DB_PASS=********* --set-env-vars INSTANCE_CONNECTION_NAME=***********:us-central1:topshot-****** --set-env-vars DB_NAME=topshot******** --set-env-vars DEFAULT_BLOCK_HEIGHT_STREAM_SIZE=200  --set-env-vars DEFAULT_BLOCK_HEIGHT_PLAYSET_QUERY=60000
```

This will allow public access to the API.  Ideally, you could put this behind an API gateway and enable user authentication.

