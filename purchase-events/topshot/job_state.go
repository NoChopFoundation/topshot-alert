/**
 * author: NoChopFoundation@gmail.com
 */
package topshot

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const WORK_AREA_FOLDER = "state"

var JOB_STATE_FILE_PATH = filepath.Join(WORK_AREA_FOLDER, "jobState.json")

type JobState struct {
	LastBlockProcessed uint64
}

func Initialize() {
	err := os.MkdirAll(WORK_AREA_FOLDER, 0700)
	if err != nil {
		panic(err)
	}

	log.SetFormatter(&log.JSONFormatter{})
}

func LoadJobState() JobState {
	var prevJobstate JobState
	if _, err := os.Stat(JOB_STATE_FILE_PATH); err == nil {
		jsonFile, err := os.Open(JOB_STATE_FILE_PATH)
		if err != nil {
			panic(err)
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal([]byte(byteValue), &prevJobstate)
	} else {
		prevJobstate.LastBlockProcessed = 0
	}
	return prevJobstate
}

func SaveJobState(currentJobState *JobState) {
	file, _ := json.MarshalIndent(currentJobState, "", " ")
	err := ioutil.WriteFile(JOB_STATE_FILE_PATH, file, 0644)
	if err != nil {
		panic(err)
	}
}
