# TopShot Purchase/Listing Event Collector
This program queries the Flow blockchain for events related to the NBA TopShots contract.  This includes events when a moment is listed or purchased.


## Purpose
To monitor and gain insight on the price of NBA TopShot moments.  By default it simply prints out recent events to the console but can also be configured to write events to a MySQL database.

## Caveats
- This was my first time working with the Go programming language so it might be considered sloppy Go code.
- I used the Flow SDK to query events.  I do not think they support getting historical event data using the API ExecuteScriptAtBlockHeight.  My understanding is the Flow query servers only keep the full state of fairly recent blocks (which is required for event generation).  Since it is a block chain there must be some way of downloading the entire data and then processing each block in order to see the old events (like minting events).  My database has been running since April 2021, so it is possible I could provide my current state.

## Usage

### Setup
- Install Go
- Execute from source
  > go run query_loop.go 
- Build an executable
  > go build query_loop.go


### Arguments

If no arguments are specified it will print recent purchase/listing events to the console.

```
  -dumpSchema
        console print database schema and latest moment data in TopShot contract  
  -i int
        the number of iterations to execute before exiting (0 forever) (default 1)
  -quiet
        do not display events to the console
  -useDb
        attempt to use a database, furthur env variables are required
```

### Display Recent Events

```
PS ..\purchase-events> go run .\query_loop.go   
{"level":"info","msg":"loop main","prevJobstate.LastBlockProcessed":14732429,"time":"2021-05-21T19:24:32-04:00"}
{"latestBlock.Height":100,"level":"info","msg":"calling query","time":"2021-05-21T19:24:32-04:00"}
{"N":"Ricky Rubio","P":2,"SN":13903,"T":"Base Set","level":"info","msg":"Purchased","time":"2021-05-21T19:24:32-04:00"}
{"N":"Jarrett Allen","P":1,"SN":22766,"T":"Base Set","level":"info","msg":"Purchased","time":"2021-05-21T19:24:32-04:00"}
{"N":"Jrue Holiday","P":2,"SN":9978,"T":"Base Set","level":"info","msg":"Purchased","time":"2021-05-21T19:24:32-04:00"}
{"N":"Karl-Anthony Towns","P":3,"SN":20242,"T":"Base Set","level":"info","msg":"Purchased","time":"2021-05-21T19:24:32-04:00"}
{"N":"Victor Oladipo","P":2,"SN":33913,"T":"Base Set","level":"info","msg":"Purchased","time":"2021-05-21T19:24:33-04:00"}
{"N":"Ben Simmons","P":19,"SN":9458,"T":"Cool Cats","level":"info","msg":"Purchased","time":"2021-05-21T19:24:33-04:00"}
{"N":"Tyler Herro","P":9,"SN":13344,"T":"Base Set","level":"info","msg":"Purchased","time":"2021-05-21T19:24:33-04:00"}
{"N":"Matisse Thybulle","P":20,"SN":2479,"T":"Hustle and Show","level":"info","msg":"Purchased","time":"2021-05-21T19:24:33-04:00"}
{"N":"Jayson Tatum","P":2,"SN":34081,"T":"Base Set","level":"info","msg":"Purchased","time":"2021-05-21T19:24:33-04:00"}
{"N":"Tacko Fall","P":13,"SN":4988,"T":"Base Set","level":"info","msg":"Purchased","time":"2021-05-21T19:24:33-04:00"}
{"N":"Kristaps Porziņģis","P":35,"SN":3326,"T":"Base Set","level":"info","msg":"Listed","time":"2021-05-21T19:24:33-04:00"}
{"N":"Karl-Anthony Towns","P":21,"SN":17612,"T":"Hustle and Show","level":"info","msg":"Listed","time":"2021-05-21T19:24:33-04:00"}
```

### Enumerate Existing Moments
The moment metadata can be dumped via ***-dumpSchema*** option.

```
PS ..\purchase-events> go run .\query_loop.go -dumpSchema > dump.txt
```
dump.txt
```
INSERT into `plays` (`PlayId`, `FullName`, `FirstName`, `LastName`, `Birthdate`, `Birthplace`, `JerseyNumber`, `DraftTeam`, `DraftYear`, `DraftSelection`, `DraftRound`, `TeamAtMomentNBAID`, `TeamAtMoment`, `PrimaryPosition`, `PlayerPosition`, `Height`, `Weight`, `TotalYearsExperience`, `NbaSeason`, `DateOfMoment`, `PlayCategory`, `PlayType`, `HomeTeamName`, `AwayTeamName`, `HomeTeamScore`, `AwayTeamScore`) 
VALUES (1, 'Trae Young', 'Trae', 'Young', '1998-09-19', 'Lubbock, TX, USA', '11', 'Dallas Mavericks', '2018', '5', '1', '1610612737', 'Atlanta Hawks', 'PG', 'G', '73', '180', '1', '2019-20', '2019-11-06 00:30:00 +0000 UTC', 'Handles', 'Handles', 'Atlanta Hawks', 'San Antonio Spurs', '108', '100') ON DUPLICATE KEY UPDATE created_at=NOW(); 

INSERT into `plays` (`PlayId`, `FullName`, `FirstName`, `LastName`, `Birthdate`, `Birthplace`, `JerseyNumber`, `DraftTeam`, `DraftYear`, `DraftSelection`, `DraftRound`, `TeamAtMomentNBAID`, `TeamAtMoment`, `PrimaryPosition`, `PlayerPosition`, `Height`, `Weight`, `TotalYearsExperience`, `NbaSeason`, `DateOfMoment`, `PlayCategory`, `HomeTeamName`, `AwayTeamName`, `HomeTeamScore`, `AwayTeamScore`, `PlayType`) 
VALUES (32, 'Maxi Kleber', 'Maxi', 'Kleber', '1992-01-29', 'Wurzburg,, DEU', '42', 'N/A', '2014', 'N/A', 'N/A', '1610612742', 'Dallas Mavericks', 'C', 'F', '82', '240', '2', '2019-20', '2019-12-04 00:30:00 +0000 UTC', 'Dunk', 'New Orleans Pelicans', 'Dallas Mavericks', '97', '118', 'Rim') ON DUPLICATE KEY UPDATE created_at=NOW(); 

....

INSERT into sets (`SetId`, `Name`, `SetSeries`) values (1, 'Genesis', 0) ON DUPLICATE KEY UPDATE created_at=NOW();
INSERT into sets (`SetId`, `Name`, `SetSeries`) values (2, 'Base Set', 0) ON DUPLICATE KEY UPDATE created_at=NOW();

.....

INSERT into plays_in_sets (`SetId`, `PlayId`, `EditionCount`) values (1, 32, 1) ON DUPLICATE KEY UPDATE created_at=NOW();
INSERT into plays_in_sets (`SetId`, `PlayId`, `EditionCount`) values (1, 33, 1) ON DUPLICATE KEY UPDATE created_at=NOW();

```

## Database Population
It is also possible to populate a database with the emitted purchase/listing events.

General Steps (tested with MySQL)
- Install MySQL
- Create the event tables (see ***sql/event_tables_mysql.sql***)
- Set environment variable DB_CONNECTION
  - For local dev, DB_CONNECTION = root:***@tcp(127.0.0.1:3306)/testdb
- Set environment variable COLLECTOR_ID
  - This is any user defined token which identifies the current machine, COLLECTOR_ID = 1.
  - Every few minutues it will update it's status to the database so it is easy to test if the collector is up.
- Execute with argument ***-useDb*** and ***-i 0*** to execute forever.
  - In practice, I do not know the memory leak properties so I just used ***-i 200*** and wrapped the execution in a shell loop.


## Credits
This was inspired and contains elements from https://medium.com/@eric.ren_51534/polling-nba-top-shot-p2p-market-purchase-events-from-flow-blockchain-using-flow-go-sdk-3ec80119e75f
