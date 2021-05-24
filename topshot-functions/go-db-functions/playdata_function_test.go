package playdata

// Must have dev DB loaded

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
)

// TODO This testcase requires the collector to be running as we test
//      if any collectors have reported in the last minute
func TestGetPlayDataHTTP_RecentCollectors(t *testing.T) {
	os.Setenv("DB_CONNECTION", "root:*******@tcp(127.0.0.1:3306)/testdb")
	_, got := get("/status/collectors/recent")

	var payload PayloadCollector
	toJson(t, got, &payload)

	if len(payload.Data) < 1 {
		t.Error("Not enough data in db")
	}
	for _, event := range payload.Data {
		assertNotZero64(t, event.BlockHeight)
		assetNotEmpty(t, event.State)
	}
}

func TestGetPlayDataHTTP_RecentByPlaySet(t *testing.T) {
	os.Setenv("DB_CONNECTION", "root:******@tcp(127.0.0.1:3306)/testdb")
	os.Setenv("DEFAULT_BLOCK_HEIGHT_STREAM_SIZE", "1000")
	_, got := get("/stream/momentEvents/from/HEAD")

	var payload PayloadMomentEvents
	toJson(t, got, &payload)

	if len(payload.Data) < 5 {
		t.Error("Not enough data in db")
	}

	setId := payload.Data[0].SetId
	playId := payload.Data[0].PlayId
	_, got2 := get("/momentEvents/play/" + strconv.Itoa(playId) +
		"/set/" + strconv.Itoa(setId) + "/from/HEAD")
	toJson(t, got2, &payload)
	if len(payload.Data) < 1 {
		t.Error("Not enough data in db")
	}

	for _, event := range payload.Data {
		assertNotZero(t, event.MomentId)
		assetNotEmpty(t, event.Type)
		assetNotEmpty(t, event.SellerAddr)
		assertNotZero64(t, event.BlockHeight)
		assertNotZero(t, event.PlayId)
		assertNotZero(t, event.SetId)
		assertNotZero(t, event.SerialNumber)
		assertNotZeroFloat(t, event.Price)
		assetNotEmpty(t, event.Created_At)

		assetEqual(t, event.PlayId, playId)
		assetEqual(t, event.SetId, setId)
	}
}

func TestGetPlayDataHTTP_StreamMaxBlockHeight(t *testing.T) {
	os.Setenv("DB_CONNECTION", "root:******@tcp(127.0.0.1:3306)/testdb")
	os.Setenv("DEFAULT_BLOCK_HEIGHT_STREAM_SIZE", "1000")
	_, got := get("/stream/momentEvents/from/HEAD")

	var payload PayloadMomentEvents
	toJson(t, got, &payload)

	if len(payload.Data) < 5 {
		t.Error("Not enough data in db")
	}

	for _, event := range payload.Data {
		assertNotZero(t, event.MomentId)
		assetNotEmpty(t, event.Type)
		assetNotEmpty(t, event.SellerAddr)
		assertNotZero64(t, event.BlockHeight)
		assertNotZero(t, event.PlayId)
		assertNotZero(t, event.SetId)
		assertNotZero(t, event.SerialNumber)
		assertNotZeroFloat(t, event.Price)
		assetNotEmpty(t, event.Created_At)
	}
	i := 0
	for i < 100 { // Cause some stress on DB pool logic
		_, got2 := get("/stream/momentEvents/from/HEAD/to/blockHeight/" + strconv.Itoa(int(payload.Data[0].BlockHeight-1)))

		var payload2 PayloadMomentEvents
		toJson(t, got2, &payload2)

		if len(payload2.Data) < 1 {
			t.Error("Not enough data in db")
		}

		if len(payload.Data) <= len(payload2.Data) {
			t.Error("expecting targeted request to be much smaller")
		}

		for _, event := range payload2.Data {
			assertNotZero(t, event.MomentId)
			assetNotEmpty(t, event.Type)
			assetNotEmpty(t, event.SellerAddr)
			assertNotZero64(t, event.BlockHeight)
			assertNotZero(t, event.PlayId)
			assertNotZero(t, event.SetId)
			assertNotZero(t, event.SerialNumber)
		}

		i = i + 1
	}
}

func TestGetPlayDataHTTP_Plays(t *testing.T) {
	os.Setenv("DB_CONNECTION", "root:******@tcp(127.0.0.1:3306)/testdb")
	_, got := get("/plays")
	_, got2 := get("/plays")

	var payload, cachedPayload PayloadPlays
	toJson(t, got, &payload)
	toJson(t, got2, &cachedPayload)

	if len(payload.Data) < 100 {
		t.Error("Not enough data in db")
	}
	if len(payload.Data) != len(cachedPayload.Data) {
		t.Error("cache not working")
	}

	for _, play := range payload.Data {
		if play.PlayId == 0 {
			t.Error("PlayId not set")
		}
		assetNotEmpty(t, play.FullName)
		assetNotEmpty(t, play.TeamAtMoment)

		_, gotPlay := get("/plays/" + strconv.Itoa(play.PlayId))
		var currentPlay Play
		toJson(t, gotPlay, &currentPlay)

		if play.PlayId != currentPlay.PlayId {
			t.Error("PlayId not equal")
		}
		if play.FullName != currentPlay.FullName {
			t.Error("Set not equal")
		}
		if len(play.Sets) == 0 {
			t.Error("No sets data")
		}
		if len(play.EditionCounts) != len(play.Sets) {
			t.Error("No edition data")
		}
	}
}

func TestGetPlayDataHTTP_Sets(t *testing.T) {
	os.Setenv("DB_CONNECTION", "root:******@tcp(127.0.0.1:3306)/testdb")
	_, got := get("/sets")

	var payload PayloadSets
	toJson(t, got, &payload)

	if len(payload.Data) < 10 {
		t.Error("Not enough data in db")
	}

	for _, set := range payload.Data {
		if set.SetId == 0 {
			t.Error("SetId not set")
		}
		assetNotEmpty(t, set.Name)
		assertSetSeries(t, set.SetId, set.SetSeries)

		_, gotSet := get("/sets/" + strconv.Itoa(set.SetId))
		var currentSet Set
		toJson(t, gotSet, &currentSet)

		if set.SetId != currentSet.SetId {
			t.Error("SetId not equal")
		}
		if set.Name != currentSet.Name {
			t.Error("Set not equal")
		}

		_, gotSetPlayers := get("/sets/" + strconv.Itoa(set.SetId) + "/plays")
		var playsPayload PayloadSetWithPlay
		toJson(t, gotSetPlayers, &playsPayload)
		if len(playsPayload.Data) < 2 {
			t.Error("Not enough play in set data in db")
		}
		for _, playInset := range playsPayload.Data {
			assetNotEmpty(t, playInset.FullName)
			if playInset.PlayId == 0 {
				t.Error("set plat id not set")
			}
		}
	}
}

/**
|     1 | Genesis                        |         0 | 2021-04-03 14:21:29 |
|     2 | Base Set                       |         0 | 2021-04-03 14:21:29 |
|     3 | Platinum Ice                   |         0 | 2021-04-03 14:21:29 |
|     4 | Holo MMXX                      |         0 | 2021-04-03 14:21:29 |
|     5 | Metallic Gold LE               |         0 | 2021-04-03 14:21:29 |
|     6 | Early Adopters                 |         0 | 2021-04-03 14:21:29 |
|     7 | Rookie Debut                   |         0 | 2021-04-03 14:21:29 |
|     8 | Cosmic                         |         0 | 2021-04-03 14:21:29 |
|     9 | For the Win                    |         0 | 2021-04-03 14:21:30 |
|    10 | Denied!                        |         0 | 2021-04-03 14:21:30 |
|    11 | Throwdowns                     |         0 | 2021-04-03 14:21:30 |
|    12 | From the Top                   |         0 | 2021-04-03 14:21:30 |
|    13 | With the Strip                 |         0 | 2021-04-03 14:21:30 |
|    14 | Hometown Showdown: Cali vs. NY |         0 | 2021-04-03 14:21:30 |
|    15 | So Fresh                       |         0 | 2021-04-03 14:21:30 |
|    16 | First Round                    |         0 | 2021-04-03 14:21:30 |
|    17 | Conference Semifinals          |         0 | 2021-04-03 14:21:30 |
|    18 | Western Conference Finals      |         0 | 2021-04-03 14:21:30 |
|    19 | Eastern Conference Finals      |         0 | 2021-04-03 14:21:30 |
|    20 | 2020 NBA Finals                |         0 | 2021-04-03 14:21:30 |
|    21 | The Finals                     |         0 | 2021-04-03 14:21:30 |
|    22 | Got Game                       |         0 | 2021-04-03 14:21:30 |
|    23 | Lace 'Em Up                    |         0 | 2021-04-03 14:21:30 |
|    24 | MVP Moves                      |         0 | 2021-04-03 14:21:30 |
|    25 | Run It Back                    |         0 | 2021-04-03 14:21:30 |
|    26 | Base Set                       |         2 | 2021-04-03 14:21:30 |
|    27 | Platinum Ice                   |         2 | 2021-04-03 14:21:30 |
|    28 | Holo Icon                      |         2 | 2021-04-03 14:21:30 |
|    29 | Metallic Gold LE               |         2 | 2021-04-03 14:21:30 |
|    30 | Season Tip-off                 |         2 | 2021-04-03 14:21:30 |
|    31 | Deck the Hoops                 |         2 | 2021-04-03 14:21:30 |
|    32 | Cool Cats                      |         2 | 2021-04-03 14:21:30 |
|    33 | The Gift                       |         2 | 2021-04-03 14:21:30 |
|    34 | Seeing Stars                   |         2 | 2021-04-03 14:21:30 |
|    35 | Rising Stars                   |         2 | 2021-04-03 14:21:30 |
|    36 | 2021 All-Star Game             |         2 | 2021-04-03 14:21:30 |
*/

func assertSetSeries(t *testing.T, setId int, setSeries int) {
	if setId < 26 {
		if setSeries != 0 {
			t.Errorf("assertSetSeries %d %d", setId, setSeries)
		}
	} else if setId < 36 {
		if setSeries != 2 {
			t.Errorf("assertSetSeries %d %d", setId, setSeries)
		}
	}
}

func assetEqual(t *testing.T, a int, b int) {
	if a != b {
		t.Errorf("actual %d %d", a, b)
	}
}

func assetNotEmpty(t *testing.T, str string) {
	if len(str) == 0 {
		t.Errorf("actual %s", str)
	}
}

func assertNotZero(t *testing.T, v int) {
	if v == 0 {
		t.Errorf("actual %d", v)
	}
}

func assertNotZeroFloat(t *testing.T, v float32) {
	if v == 0 {
		t.Errorf("actual %f", v)
	}
}

func assertNotZero64(t *testing.T, v uint64) {
	if v == 0 {
		t.Errorf("actual %d", v)
	}
}
func get(path string) (*httptest.ResponseRecorder, string) {
	req := httptest.NewRequest("GET", path, strings.NewReader(""))
	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	GetPlayDataHTTP(rr, req)

	return rr, rr.Body.String()
}

func toJson(t *testing.T, got string, into interface{}) {
	err := json.Unmarshal([]byte(got), into)
	if err != nil {
		t.Error("Unable to convert to json")
	}
}
