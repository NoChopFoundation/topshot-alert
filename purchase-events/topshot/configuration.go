/**
 * author: NoChopFoundation@gmail.com
 */
package topshot

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Configuration struct {
	// access.mainnet.nodes.onflow.org:9000 etc
	accessNode string

	// root:*****@tcp(127.0.0.1:3306)/testdb
	MySqlConnection string

	CollectorId uint16
}

func Configuration_MainNet() (*Configuration, error) {
	config := Configuration{
		accessNode: "access.mainnet.nodes.onflow.org:9000",
	}
	return &config, nil
}

func Configuration_MainNet_withMySql() (*Configuration, error) {
	dbConnStr, err := getMySqlConnectionStr()
	if err != nil {
		return nil, err
	}

	collectorId, err := strconv.Atoi(os.Getenv("COLLECTOR_ID"))
	if err != nil {
		return nil, errors.New("require COLLECTOR_ID unique int for collector status updates")
	}

	config := Configuration{
		accessNode:      "access.mainnet.nodes.onflow.org:9000",
		MySqlConnection: dbConnStr,
		CollectorId:     uint16(collectorId),
	}
	return &config, nil
}

func Configuration_ConnectionError() (*Configuration, error) {
	config := Configuration{
		accessNode: "access.mainnet.nodes.onflow.org:8999",
	}
	return &config, nil
}

// If DB_CONNECTION is set, we assume it is from local dev environment, like "root:***@tcp(127.0.0.1:3306)/testdb"
// If INSTANCE_CONNECTION_NAME is set, then we assume its gcloud which also sets
//   DB_USER, DB_PASS, DB_NAME where INSTANCE_CONNECTION_NAME like primeval-pen-307423:us-central1:topshot2
func getMySqlConnectionStr() (string, error) {
	dbConnectionStr := os.Getenv("DB_CONNECTION")
	if len(dbConnectionStr) == 0 {
		var (
			dbUser                 = os.Getenv("DB_USER")                  // e.g. 'my-db-user'
			dbPwd                  = os.Getenv("DB_PASS")                  // e.g. 'my-db-password'
			instanceConnectionName = os.Getenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
			dbName                 = os.Getenv("DB_NAME")                  // e.g. 'my-database'
		)
		if len(instanceConnectionName) == 0 {
			return "", errors.New("require environment variable DB_CONNECTION or (DB_USER,DB_PASS,INSTANCE_CONNECTION_NAME,DB_NAME)")
		}
		// GCloud deployment
		socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
		if !isSet {
			socketDir = "/cloudsql"
		}

		dbConnectionStr = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", dbUser, dbPwd, socketDir, instanceConnectionName, dbName)
		return dbConnectionStr, nil
	} else {
		// Local development use-case
		return dbConnectionStr, nil
	}
}
