package playdata

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

// init runs during package initialization. So, this will only run during an
// an instance's cold start.
func init() {
	Cache = new(GlobalCache)
	Cache.mutex = &sync.Mutex{}
}

func GetPlayDataHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO Cors in development mode, allow every from where
	// Set CORS headers for the preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		// Set CORS headers for the main request.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		doGetPlayDataHTTP(w, r)
	}
}

func doGetPlayDataHTTP(w http.ResponseWriter, r *http.Request) {
	router := httprouter.New()
	router.GET("/plays", func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		cachedPlays := GlobalCacheAvailable(Cache)
		if cachedPlays != nil {
			handleSendJson(cachedPlays, w)
		} else {
			databaseHandler(handleGetAllPlays, w, r, ps)
		}
	})
	router.GET("/plays/:PlayId", func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		databaseHandler(handleGetOnePlay, w, r, ps)
	})
	router.GET("/sets", func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		databaseHandler(handleGetAllSets, w, r, ps)
	})
	router.GET("/sets/:SetId", func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		databaseHandler(handleGetOneSet, w, r, ps)
	})
	router.GET("/sets/:SetId/plays", func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		databaseHandler(handleGetOneSetWithPlayerData, w, r, ps)
	})
	router.GET("/stream/momentEvents/from/HEAD", func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		databaseHandler(handleGetStreamFromMaxBlockHeight, w, r, ps)
	})
	router.GET("/stream/momentEvents/from/HEAD/to/blockHeight/:BlockHeight", func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		databaseHandler(handleGetStreamFromMaxBlockHeight, w, r, ps)
	})
	router.GET("/status/collectors/recent", func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		databaseHandler(handleGetRecentCollectors, w, r, ps)
	})
	router.GET("/momentEvents/play/:PlayId/set/:SetId/from/HEAD", func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		databaseHandler(handleGetRecentPlaySet, w, r, ps)
	})

	router.ServeHTTP(w, r)
}

func databaseHandler(callBack DbHandlerFunc, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	opts, db, err := getDbConnectionInit()
	if err != nil {
		handleInternalError("db init", err, w)
	} else {
		callBack(DbHandlerParams{
			Opts:       opts,
			Db:         db,
			W:          w,
			Request:    r,
			RouteParms: ps})
	}
}

func GetSet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "GetPlay, %s!\n", ps.ByName("SetId"))
}

func getDbConnectionInit() (*GetPlayDataHTTPInit, *sql.DB, error) {
	opts, err := getInit()
	if err != nil {
		return nil, nil, err
	} else {
		var dbPool *sql.DB
		dbFromCache := GlobalCache_GetDb(Cache)
		if dbFromCache != nil {
			dbPool = dbFromCache
		} else {
			dbPool, err = sql.Open("mysql", opts.dbConnectionStr)
			if err != nil {
				return nil, nil, err
			} else {
				dbPool.SetMaxOpenConns(5)
				dbPool.SetConnMaxLifetime(1800)
				dbPool.SetMaxIdleConns(3)
				GlobalCache_SetDb(Cache, dbPool)
			}
		}
		return &opts, dbPool, nil
	}
}

// 1st grab which plays belong in which sets
// 2nd grab plays metadata
// 3rd tack on set data
func handleGetAllPlays(p DbHandlerParams) {
	queryHandler(QueryHandlerParms{
		query:     GetPlaysInSets(p.Opts),
		queryArgs: nil,
		p:         p,
		callback: func(rows *sql.Rows, p DbHandlerParams) {
			playsInSets, err := scanPlaysInSetsRows(rows)
			if err != nil {
				handleInternalError("Internal error query select row", err, p.W)
			} else {
				queryHandler(QueryHandlerParms{
					query:     GetPlaySelect(p.Opts),
					queryArgs: nil,
					p:         p,
					callback: func(rows *sql.Rows, p DbHandlerParams) {
						payload, err := scanPlayRows(rows)
						if err != nil {
							handleInternalError("Internal error query select row", err, p.W)
						} else {
							for i, _ := range payload.Data {
								payload.Data[i].Sets = tackOnSetData(payload.Data[i].PlayId, playsInSets)
								payload.Data[i].EditionCounts = tackOnEditionData(payload.Data[i].PlayId, playsInSets)
							}
							GlobalCacheSet(Cache, payload)
							handleSendJson(payload, p.W)
						}
					},
				})
			}
		},
	})
}

func tackOnSetData(playId int, playsInSets *PayLoadPlaysInSets) []int {
	sets := []int{}
	for _, playsinSet := range playsInSets.Data {
		if playsinSet.PlayId == playId {
			sets = append(sets, playsinSet.SetId)
		}
	}
	return sets
}

func tackOnEditionData(playId int, playsInSets *PayLoadPlaysInSets) []int {
	editionCounts := []int{}
	for _, playsinSet := range playsInSets.Data {
		if playsinSet.PlayId == playId {
			editionCounts = append(editionCounts, playsinSet.EditionCount)
		}
	}
	return editionCounts
}

func handleGetAllSets(p DbHandlerParams) {
	queryHandler(QueryHandlerParms{
		query:     GetSetSelect(p.Opts),
		queryArgs: nil,
		p:         p,
		callback: func(rows *sql.Rows, p DbHandlerParams) {
			payload, err := scanSetRows(rows)
			if err != nil {
				handleInternalError("Internal error query select row", err, p.W)
			} else {
				handleSendJson(payload, p.W)
			}
		},
	})
}

func handleGetOnePlay(p DbHandlerParams) {
	playId, err := strconv.Atoi(p.RouteParms.ByName("PlayId"))
	if err != nil {
		handleBadRequest("PlayId parameter", p.W)
	} else {
		queryHandler(QueryHandlerParms{
			query:     GetPlaySelect(p.Opts) + " WHERE PlayId=?",
			queryArgs: []interface{}{playId},
			p:         p,
			callback: func(rows *sql.Rows, p DbHandlerParams) {
				payload, err := scanPlayRows(rows)
				if err != nil {
					handleInternalError("Internal error query select row", err, p.W)
				} else {
					if len(payload.Data) == 1 {
						handleSendJson(payload.Data[0], p.W)
					} else {
						handleNotFound("play doesn't exist", p.W)
					}
				}
			}})
	}
}

func handleGetOneSet(p DbHandlerParams) {
	setId, err := strconv.Atoi(p.RouteParms.ByName("SetId"))
	if err != nil {
		handleBadRequest("SetId parameter", p.W)
	} else {
		queryHandler(QueryHandlerParms{
			query:     GetSetSelect(p.Opts) + " WHERE SetId=?",
			queryArgs: []interface{}{setId},
			p:         p,
			callback: func(rows *sql.Rows, p DbHandlerParams) {
				payload, err := scanSetRows(rows)
				if err != nil {
					handleInternalError("Internal error query select row", err, p.W)
				} else {
					if len(payload.Data) == 1 {
						handleSendJson(payload.Data[0], p.W)
					} else {
						handleNotFound("set doesn't exist", p.W)
					}
				}
			}})
	}
}

func handleGetRecentCollectors(p DbHandlerParams) {
	queryHandler(QueryHandlerParms{
		query:     GetRecentCollectorStatus(p.Opts),
		queryArgs: []interface{}{p.Opts.defaultCollectorPastSeconds},
		p:         p,
		callback: func(rows *sql.Rows, p DbHandlerParams) {
			payload, err := scanCollectorRows(rows)
			if err != nil {
				handleInternalError("Internal error query select row", err, p.W)
			} else {
				handleSendJson(payload, p.W)
			}
		},
	})
}

func handleGetRecentPlaySet(p DbHandlerParams) {
	playId, err := strconv.ParseUint(p.RouteParms.ByName("PlayId"), 10, 64)
	setId, err2 := strconv.ParseUint(p.RouteParms.ByName("SetId"), 10, 64)

	if err != nil || err2 != nil {
		handleBadRequest("PlayId or SetId parameter", p.W)
	} else {
		queryHandler(QueryHandlerParms{
			query:     GetRecentByPlaySet(p.Opts),
			queryArgs: []interface{}{playId, setId, p.Opts.defaultBlockHeightPlaySetQuery},
			p:         p,
			callback: func(rows *sql.Rows, p DbHandlerParams) {
				payload, err := scanMomentEventsRows(rows)
				if err != nil {
					handleInternalError("Internal error query select row", err, p.W)
				} else {
					handleSendJson(payload, p.W)
				}
			},
		})
	}
}

func handleGetStreamFromMaxBlockHeight(p DbHandlerParams) {
	strFromBlockHeight := p.RouteParms.ByName("BlockHeight")
	if len(strFromBlockHeight) == 0 {
		// (HEAD - Small_value) > DEFAULT_BLOCK_HEIGHT_STREAM_SIZE, so we end up streaming HEAD to DEFAULT_BLOCK_HEIGHT_STREAM_SIZE
		strFromBlockHeight = "1"
	}

	fromBlockHeight, err := strconv.ParseUint(strFromBlockHeight, 10, 64)
	if err != nil {
		handleBadRequest("BlockHeight parameter", p.W)
	} else {
		queryHandler(QueryHandlerParms{
			query:     GetRecentMoments(p.Opts),
			queryArgs: []interface{}{fromBlockHeight, p.Opts.defaultBlockHeightStreamSize},
			p:         p,
			callback: func(rows *sql.Rows, p DbHandlerParams) {
				payload, err := scanMomentEventsRows(rows)
				if err != nil {
					handleInternalError("Internal error query select row", err, p.W)
				} else {
					handleSendJson(payload, p.W)
				}
			},
		})
	}
}

func handleGetOneSetWithPlayerData(p DbHandlerParams) {
	setId, err := strconv.Atoi(p.RouteParms.ByName("SetId"))
	if err != nil {
		handleBadRequest("SetId parameter", p.W)
	} else {
		queryHandler(QueryHandlerParms{
			query:     GetSetSelectWithPlayerData(p.Opts) + " WHERE " + DBTABLE_plays_in_sets + ".SetId=?",
			queryArgs: []interface{}{setId},
			p:         p,
			callback: func(rows *sql.Rows, p DbHandlerParams) {
				payload, err := scanSetPlayRows(rows)
				if err != nil {
					handleInternalError("Internal error query select row", err, p.W)
				} else {
					if len(payload.Data) > 0 {
						handleSendJson(payload, p.W)
					} else {
						handleNotFound("set doesn't exist", p.W)
					}
				}
			}})
	}
}

type QueryHandlerFunc func(rows *sql.Rows, p DbHandlerParams)

type QueryHandlerParms struct {
	query     string
	queryArgs []interface{}
	p         DbHandlerParams
	callback  QueryHandlerFunc
}

func queryHandler(args QueryHandlerParms) {
	var rows *sql.Rows
	var err error
	if args.queryArgs != nil {
		rows, err = args.p.Db.Query(args.query, args.queryArgs...)
	} else {
		rows, err = args.p.Db.Query(args.query)
	}

	if err != nil {
		handleInternalError("Internal error query select", err, args.p.W)
	} else {
		args.callback(rows, args.p)
		rows.Close()
	}
}

func scanSetPlayRows(rows *sql.Rows) (*PayloadSetWithPlay, error) {
	payload := PayloadSetWithPlay{
		Data: []SetPlay{},
	}
	for rows.Next() {
		var s SetPlay
		err := rows.Scan(&s.PlayId, &s.FullName)
		if err != nil {
			return nil, err
		}
		payload.Data = append(payload.Data, s)
	}
	return &payload, nil
}

func scanSetRows(rows *sql.Rows) (*PayloadSets, error) {
	payload := PayloadSets{
		Data: []Set{},
	}
	for rows.Next() {
		var s Set
		err := rows.Scan(&s.SetId, &s.Name, &s.SetSeries)
		if err != nil {
			return nil, err
		}
		payload.Data = append(payload.Data, s)
	}
	return &payload, nil
}

func scanPlaysInSetsRows(rows *sql.Rows) (*PayLoadPlaysInSets, error) {
	payload := PayLoadPlaysInSets{
		Data: []PlaysInSets{},
	}
	for rows.Next() {
		var p PlaysInSets
		err := rows.Scan(&p.PlayId, &p.SetId, &p.EditionCount)
		if err != nil {
			return nil, err
		}
		payload.Data = append(payload.Data, p)
	}
	return &payload, nil
}

func scanCollectorRows(rows *sql.Rows) (*PayloadCollector, error) {
	payload := PayloadCollector{
		Data: []Collector{},
	}
	for rows.Next() {
		var p Collector
		err := rows.Scan(&p.CollectorId, &p.State, &p.UpdatesInInterval, &p.BlockHeight, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		payload.Data = append(payload.Data, p)
	}
	return &payload, nil
}

func scanMomentEventsRows(rows *sql.Rows) (*PayloadMomentEvents, error) {
	payload := PayloadMomentEvents{
		Data: []MomentEvents{},
	}
	for rows.Next() {
		var p MomentEvents
		err := rows.Scan(&p.Type, &p.MomentId, &p.BlockHeight, &p.PlayId,
			&p.SerialNumber, &p.SetId, &p.SellerAddr, &p.Price, &p.Created_At)
		if err != nil {
			return nil, err
		}
		payload.Data = append(payload.Data, p)
	}
	return &payload, nil
}

func scanPlayRows(rows *sql.Rows) (*PayloadPlays, error) {
	payload := PayloadPlays{
		Data: []Play{},
	}
	for rows.Next() {
		var p PlayScan
		err := rows.Scan(&p.PlayId, &p.NbaSeason, &p.TeamAtMomentNBAID, &p.PlayCategory,
			&p.JerseyNumber, &p.PlayerPosition, &p.DateOfMoment,
			&p.PlayType, &p.FullName, &p.PrimaryPosition, &p.TeamAtMoment)
		if err != nil {
			return nil, err
		}
		payload.Data = append(payload.Data, removeNulls(&p))
	}
	return &payload, nil
}

func handleSendJson(payload interface{}, w http.ResponseWriter) {
	resp, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal error json marshal", http.StatusInternalServerError)
	} else {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, string(resp))
	}
}

func removeNulls(pScan *PlayScan) Play {
	// Scan doesn't like null columns, should do maps
	return Play{
		PlayId:            pScan.PlayId,
		NbaSeason:         pScan.NbaSeason.String,
		TeamAtMomentNBAID: pScan.TeamAtMomentNBAID.String,
		PlayCategory:      pScan.PlayCategory.String,
		JerseyNumber:      pScan.JerseyNumber.String,
		PlayerPosition:    pScan.PlayerPosition.String,
		DateOfMoment:      pScan.DateOfMoment.String,
		PlayType:          pScan.PlayType.String,
		FullName:          pScan.FullName.String,
		PrimaryPosition:   pScan.PrimaryPosition.String,
		TeamAtMoment:      pScan.TeamAtMoment.String,
	}
}

// If DB_CONNECTION is set, we assume it is from local dev environment, like "root:***@tcp(127.0.0.1:3306)/testdb"
// If INSTANCE_CONNECTION_NAME is set, then we assume its gcloud which also sets
//   DB_USER, DB_PASS, DB_NAME where INSTANCE_CONNECTION_NAME like primeval-pen-307423:us-central1:topshot2
func getInit() (GetPlayDataHTTPInit, error) {
	defaultBlockHeightStreamSize, err := getOSEnvInt("DEFAULT_BLOCK_HEIGHT_STREAM_SIZE", "1")
	if err != nil {
		return GetPlayDataHTTPInit{}, errors.New("require environment variable DEFAULT_BLOCK_HEIGHT_STREAM_SIZE")
	}

	defaultCollectorPastSeconds, err := getOSEnvInt("DEFAULT_COLLECTORS_PAST_SECONDS", "60")
	if err != nil {
		return GetPlayDataHTTPInit{}, errors.New("require environment variable DEFAULT_COLLECTORS_PAST_SECONDS")
	}

	defaultBlockHeightPlaySetQuery, err := getOSEnvInt("DEFAULT_BLOCK_HEIGHT_PLAYSET_QUERY", "20000")
	if err != nil {
		return GetPlayDataHTTPInit{}, errors.New("require environment variable DEFAULT_BLOCK_HEIGHT_PLAYSET_QUERY")
	}

	dbConnectionStr := os.Getenv("DB_CONNECTION")
	if len(dbConnectionStr) == 0 {
		var (
			dbUser                 = os.Getenv("DB_USER")                  // e.g. 'my-db-user'
			dbPwd                  = os.Getenv("DB_PASS")                  // e.g. 'my-db-password'
			instanceConnectionName = os.Getenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
			dbName                 = os.Getenv("DB_NAME")                  // e.g. 'my-database'
		)
		if len(instanceConnectionName) == 0 {
			return GetPlayDataHTTPInit{}, errors.New("require environment variable DB_CONNECTION or INSTANCE_CONNECTION_NAME")
		} else {
			// GCloud deployment
		}

		socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
		if !isSet {
			socketDir = "/cloudsql"
		}

		dbConnectionStr = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", dbUser, dbPwd, socketDir, instanceConnectionName, dbName)
		return GetPlayDataHTTPInit{
			defaultBlockHeightStreamSize:   defaultBlockHeightStreamSize,
			defaultCollectorPastSeconds:    defaultCollectorPastSeconds,
			defaultBlockHeightPlaySetQuery: defaultBlockHeightPlaySetQuery,
			dbConnectionStr:                dbConnectionStr}, nil
	} else {
		// Local development use-case
		return GetPlayDataHTTPInit{
			defaultBlockHeightStreamSize:   defaultBlockHeightStreamSize,
			defaultCollectorPastSeconds:    defaultCollectorPastSeconds,
			defaultBlockHeightPlaySetQuery: defaultBlockHeightPlaySetQuery,
			dbConnectionStr:                dbConnectionStr}, nil
	}
}

func getOSEnvInt(envName string, envDefault string) (int, error) {
	envValue, isSet := os.LookupEnv(envName)
	if !isSet {
		envValue = envDefault
	}
	return strconv.Atoi(envValue)
}

func handleInternalError(msg string, err error, w http.ResponseWriter) {
	log.Fatal(err)
	http.Error(w, msg, http.StatusInternalServerError)
}

func handleBadRequest(msg string, w http.ResponseWriter) NoOp {
	log.Println(msg)
	http.Error(w, msg, http.StatusBadRequest)
	return NoOp{}
}

func handleNotFound(msg string, w http.ResponseWriter) {
	log.Println(msg)
	http.Error(w, msg, http.StatusNotFound)
}
