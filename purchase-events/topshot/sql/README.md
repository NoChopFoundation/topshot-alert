# procedure.sql
This file defines the database procedure ***GetRecentMoments***.  It was tested on MySQL.

This was to simplify and make efficient a SQL query which needed to return the most recent blocks.  I found it complicated within a single SQL query to provide purchase records in the range of 
> select moment_events FROM (MAX(BlockHeight) - (MIN(MAX(BlockHeight) - RequestedFrom_BlockHeight, MAX_ALLOWED_QUERY_BACKWARDS))) to MAX(BlockHeight)

## How to load into your database
> mysql --host={mySQLHostName} --user=root --password -e "source procedure.sql" {databaseName}

# event_tables_mysql.sql
This file defines the tables representing the TopShot listing and purchase events.  It also defines the collector uploading the events to the database.
