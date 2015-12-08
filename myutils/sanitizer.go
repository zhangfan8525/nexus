package myutils

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"pudding/datastructs"
	"time"
)

const (
	PATCH_INVALID    = 0x01
	PATCH_FIXED      = 0x02
	VALIDATE_USELESS = 0x20
)

func testUseful(update *datastructs.Update, url string) bool {
	data, err := json.Marshal(update)
	if err != nil {
		MyLogger.Println(err.Error())
		//LogAlert(err.Error())
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		MyLogger.Println(err.Error())
		//LogAlert(err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		MyLogger.Println(err.Error())
		//LogAlert(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		MyLogger.Println(err.Error())
		//LogAlert(err.Error())
	}

	var status datastructs.Status
	json.Unmarshal(body, &status)
	if status.Status == VALIDATE_USELESS {
		return false
	} else {
		return true
	}
}

func Sanitize(db *sql.DB, version int, url string, done chan bool) {
	MyLogger.Println("********Sanitizing********")
	//LogColorClass("********Sanitizing********", CLR_PULL, CLS_PULL)

	getAllRowsStatement, err := db.Prepare(`SELECT id, content
											FROM patch`)
	if err != nil {
		MyLogger.Println(err.Error())
		//LogAlert(err.Error())
	}
	defer getAllRowsStatement.Close()

	updateStatusStatement, err := db.Prepare(`UPDATE patch 
											  SET status=?
											  WHERE id=?`)
	if err != nil {
		MyLogger.Println(err.Error())
		//LogAlert(err.Error())
	}
	defer updateStatusStatement.Close()

	//updateVersionStatement, err := db.Prepare(`UPDATE patch
	//SET version=?
	//WHERE id=?`)
	//if err != nil {
	//MyLogger.Println(err.Error())
	//LogAlert(err.Error())
	//}
	//defer updateVersionStatement.Close()

	rows, err := getAllRowsStatement.Query()
	if err != nil {
		MyLogger.Println(err.Error())
		//LogAlert(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var patch datastructs.Patch
		err = rows.Scan(&patch.Id, &patch.Patch)
		if err != nil {
			MyLogger.Println(err.Error())
			//LogAlert(err.Error())
		}

		var update datastructs.Update
		err = json.Unmarshal([]byte(patch.Patch), &update)
		if err != nil {
			MyLogger.Println(err.Error())
			//LogAlert(err.Error())
		}
		update.Id = patch.Id
		if update.End_time != "" && update.End_time <= time.Now().Format(TimeLayout) {
			_, err = updateStatusStatement.Exec(PATCH_INVALID, update.Id)
		} else if !testUseful(&update, url) {
			_, err = updateStatusStatement.Exec(PATCH_FIXED, update.Id)
		} //else {
		//_, err = updateVersionStatement.Exec(version, update.Id)
		//}
	}

	done <- true
}
