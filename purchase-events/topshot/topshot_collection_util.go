/**
 * author: NoChopFoundation@gmail.com
 */
package topshot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/onflow/cadence"
)

const PLAYS_TABLE_NAME = "plays"
const SETS_TABLE_NAME = "sets"
const PLAYS_IN_SET_TABLE_NAME = "plays_in_sets"

func TopShotUtil_dumpSetPlaysData(config *Configuration) {
	queryApi, err := Connection(context.Background(), config)
	panicOnError(err)

	script := `
		import TopShot from 0x0b2a3299cc857e29
		pub struct TopShotState {
			pub var nextPlayID: UInt32
			pub var nextSetID: UInt32
			pub var totalSupply: UInt64
			init(nextPlayID: UInt32, nextSetID: UInt32, totalSupply: UInt64) {
				self.nextPlayID = nextPlayID
				self.nextSetID = nextSetID
				self.totalSupply = totalSupply
			}
		}
		pub fun main(): TopShotState {
			return TopShotState(nextPlayID: TopShot.nextPlayID, nextSetID: TopShot.nextSetID, totalSupply: TopShot.totalSupply)
		}	
	`
	res, err := queryApi.ExecuteScriptAtLatestBlock(ExecuteScriptAtLatestBlock_Arg{
		CadenceScript: script,
		ScriptArgs:    []cadence.Value{}})
	panicOnError(err)

	currentTopShot := TopShotState(res.(cadence.Struct))

	// The same play exists across sets
	var prevPlays = make(map[uint32]string)
	var prevMetaFields = make(map[string]string)
	var setDatas []*SetData

	for setId := 1; setId < int(currentTopShot.NumberSets()); setId++ {
		setDatas = append(setDatas, dumpSetData(queryApi, setId, prevPlays, prevMetaFields))
	}

	generateSetData(setDatas)
	generateTableDDL(setDatas, prevMetaFields)
}

func generateSetData(setDatas []*SetData) {
	for _, setData := range setDatas {
		insertStatement := "INSERT into " + SETS_TABLE_NAME + " (`SetId`, `Name`, `SetSeries`) values ("
		insertStatement += strconv.Itoa(int(setData.ID())) + ", '" + escapeInsertValue(setData.Name()) + "', " +
			strconv.Itoa(int(setData.SeriesID())) + ") "
		insertStatement += "ON DUPLICATE KEY UPDATE created_at=NOW();"
		fmt.Println(insertStatement)
	}
	fmt.Println("")

	for _, setData := range setDatas {
		editionCounts := setData.EditionCounts()
		for arrIdx, playId := range setData.Plays() {
			insertStatement := "INSERT into " + PLAYS_IN_SET_TABLE_NAME + " (`SetId`, `PlayId`, `EditionCount`) values ("
			insertStatement += strconv.Itoa(int(setData.ID())) + ", " + strconv.Itoa(int(playId)) + ", "
			insertStatement += strconv.Itoa(int(editionCounts[arrIdx])) + ") "
			insertStatement += "ON DUPLICATE KEY UPDATE created_at=NOW();"
			fmt.Println(insertStatement)
		}
	}
}

func generateTableDDL(setDatas []*SetData, prevMetaFields map[string]string) {
	ddlCreate := "CREATE TABLE IF NOT EXISTS " + PLAYS_TABLE_NAME + " (\n\tPlayId INT PRIMARY KEY"
	for tableName, _ := range prevMetaFields {
		// The keys used for Play metadata look table safe (may change in future)
		ddlCreate += ",\n\t" + tableName + " VARCHAR(255)"
	}
	ddlCreate += ",\n\tcreated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP\n)  ENGINE=INNODB;"
	ddlCreate += "\n"
	ddlCreate += "CREATE TABLE IF NOT EXISTS " + SETS_TABLE_NAME + "(\n"
	ddlCreate += "\tSetId INT PRIMARY KEY,\n"
	ddlCreate += "\tName VARCHAR(255),\n"
	ddlCreate += "\tSetSeries INT NOT NULL,\n"
	ddlCreate += "\tcreated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP\n)   ENGINE=INNODB;\n"
	ddlCreate += "\n"
	ddlCreate += "CREATE TABLE " + PLAYS_IN_SET_TABLE_NAME + " ( \n"
	ddlCreate += "\tPlayId INT NOT NULL,\n"
	ddlCreate += "\tSetId INT NOT NULL,\n"
	ddlCreate += "\tEditionCount INT NOT NULL,\n"
	ddlCreate += "\tcreated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,\n"
	ddlCreate += "\tprimary key (PlayId, SetId),\n"
	ddlCreate += "\tFOREIGN KEY (PlayId) REFERENCES " + PLAYS_TABLE_NAME + "(PlayId),\n"
	ddlCreate += "\tFOREIGN KEY (SetId) REFERENCES " + SETS_TABLE_NAME + "(SetId)\n"
	ddlCreate += ") ENGINE=INNODB;"

	fmt.Println("")
	fmt.Println("TABLE DDL GENERATION")
	fmt.Println("")
	fmt.Println(ddlCreate)
}

func dumpSetData(queryApi *TopshotQueryApi, setId int, prevPlays map[uint32]string, prevMetaFields map[string]string) *SetData {
	script := `
		import TopShot from 0x0b2a3299cc857e29
		pub struct SetData {
			pub var setName: String
			pub var setID: UInt32
			pub var setPlayIds: [UInt32]
			pub var setPlayMetaData: [{String: String}]
			pub var setSeries: UInt32?
			pub var setEditionSizes: [UInt32]
			init(setName: String, setID: UInt32, setPlayIds: [UInt32], setSeries: UInt32?) {
				self.setName = setName
				self.setID = setID
				self.setPlayIds = setPlayIds
				self.setPlayMetaData = []
				self.setEditionSizes = []
				self.setSeries = setSeries
				for playID in setPlayIds {
					self.setPlayMetaData.append( TopShot.getPlayMetaData(playID: playID)! )
					self.setEditionSizes.append( 
						TopShot.getNumMomentsInEdition(setID: setID, playID: playID)! )
				}
			}
		}
		pub fun main(setID: UInt32): SetData {
			return SetData(setName: TopShot.getSetName(setID: setID)!, setID: setID, setPlayIds: TopShot.getPlaysInSet(setID: setID)!,
			setSeries: TopShot.getSetSeries(setID: setID))
		}
	`
	res, err := queryApi.ExecuteScriptAtLatestBlock(ExecuteScriptAtLatestBlock_Arg{
		CadenceScript: script,
		ScriptArgs:    []cadence.Value{cadence.UInt32(setId)}})
	panicOnError(err)

	currentSet := SetData(res.(cadence.Struct))
	generateSetSQLInserts(&currentSet, prevPlays, prevMetaFields)

	return &currentSet
}

// Generates to console
// INSERT into `plays` (`PlayId`, `FullName`, .... `AwayTeamScore`)
// VALUES (819, 'James Wiseman', 'James', 'Wiseman', ... '104');
func generateSetSQLInserts(currentSet *SetData, prevPlays map[uint32]string, prevMetaFields map[string]string) {
	metadatasArray := currentSet.MetaDatas()
	for arrIdx, playId := range currentSet.Plays() {
		// Plays can exist in multiple sets, don't generate multiple play INSERTs
		if _, ok := prevPlays[playId]; !ok {
			prevPlays[playId] = "set"

			var insertLeft = "INSERT into `" + PLAYS_TABLE_NAME + "` (`PlayId`"
			var insertRight = "VALUES (" + strconv.Itoa(int(playId))

			dict := metadatasArray[int(arrIdx)]
			for _, pair := range dict.Pairs {
				columnName := pair.Key.ToGoValue().(string)
				prevMetaFields[columnName] = "set"

				insertLeft += ", `" + columnName + "`"
				insertRight += ", '" + escapeInsertValue(pair.Value.ToGoValue().(string)) + "'"
			}
			insertLeft += ") \n"
			insertRight += ") ON DUPLICATE KEY UPDATE created_at=NOW(); \n"
			fmt.Println(insertLeft + insertRight)
		}
	}
}

// Used when generating INSERT statement value, INSERT into TABLE (`col`) values ('val')
// Assuming tame printable input, would need to use another library to handle all escape situations
func escapeInsertValue(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
// See TopShotUtil_dumpSetPlaysData
// The script returns a structure so we can see the number of sets in the TopShot contract
type TopShotState cadence.Struct

func (s TopShotState) NumberPlays() uint32 {
	return uint32(s.Fields[0].(cadence.UInt32))
}
func (s TopShotState) NumberSets() uint32 {
	return uint32(s.Fields[1].(cadence.UInt32))
}
func (s TopShotState) NumberMoments() uint64 {
	return uint64(s.Fields[2].(cadence.UInt64))
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
// See dumpSetData()
// The script returns a structure so we can get the TopShot.Plays and its metadata for a set
// Plays()[0] has metadata at MetaDatas()[0]
type SetData cadence.Struct

func (s SetData) Name() string {
	return string(s.Fields[0].(cadence.String))
}
func (s SetData) ID() uint32 {
	return uint32(s.Fields[1].(cadence.UInt32))
}
func (s SetData) SeriesID() uint32 {
	optionalSeries := (s.Fields[4]).(cadence.Optional)
	if series, ok := optionalSeries.Value.(cadence.UInt32); ok {
		return uint32(series)
	}
	return 999
}

func (s SetData) EditionCounts() []uint32 {
	retPlayIds := []uint32{}
	cadenceArray := (s.Fields[5]).(cadence.Array)
	for _, cadenceValue := range cadenceArray.Values {
		retPlayIds = append(retPlayIds, uint32(cadenceValue.(cadence.UInt32)))
	}
	return retPlayIds
}

func (s SetData) Plays() []uint32 {
	retPlayIds := []uint32{}
	cadenceArray := (s.Fields[2]).(cadence.Array)
	for _, cadenceValue := range cadenceArray.Values {
		retPlayIds = append(retPlayIds, uint32(cadenceValue.(cadence.UInt32)))
	}
	return retPlayIds
}

func (s SetData) MetaDatas() []cadence.Dictionary {
	retMetas := []cadence.Dictionary{}
	cadenceArray := s.Fields[3].(cadence.Array)
	for _, cadenceValue := range cadenceArray.Values {
		retMetas = append(retMetas, cadenceValue.(cadence.Dictionary))
	}
	return retMetas
}
