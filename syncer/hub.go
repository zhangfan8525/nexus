package syncer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"pudding/datastructs"
	"pudding/myutils"
)

func Put(rw http.ResponseWriter, req *http.Request) {
	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	myutils.MyLogger.Println("********File sync put********" + ip)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		myutils.MyLogger.Println(err.Error())
		fmt.Fprint(rw, ErrorResponse)
		return
	}
	myutils.MyLogger.Println("request: " + string(body))

	var sp datastructs.SyncPut
	err = json.Unmarshal(body, &sp)
	if err != nil {
		myutils.MyLogger.Println(err.Error())
		fmt.Fprint(rw, ErrorResponse)
		return
	}

	su, err := TheNexus.NewStorageUnitFromData(&sp)
	if err != nil {
		myutils.MyLogger.Println(err.Error())
		fmt.Fprint(rw, ErrorResponse)
		return
	}

	if !TheNexus.needFeedback {
		su.Valid = true
		TheNexus.toBeSaved <- su
	}

	fmt.Fprint(rw, OkResponse)
}

func Get(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		myutils.MyLogger.Println(err.Error())
		fmt.Fprint(rw, ErrorResponse)
		return
	}
	reqstr := string(body)

	var sgr datastructs.SyncGetReq
	err = json.Unmarshal(body, &sgr)
	if err != nil {
		myutils.MyLogger.Println(err.Error())
		fmt.Fprint(rw, ErrorResponse)
		return
	}

	c := TheNexus.LookupChannel(sgr.Channel)
	if c == nil {
		fmt.Fprintf(rw, ErrorResponse)
		return
	}
	sus := c.LookupUpdates(sgr.TimeStamp)
	if len(sus) > 0 {
		ip, _, _ := net.SplitHostPort(req.RemoteAddr)
		myutils.MyLogger.Println("********File sync get********" + ip)
		myutils.MyLogger.Println("request: " + reqstr)

		resp := datastructs.SyncGetResp{
			Status: 0,
		}
		for _, su := range sus {
			resp.Data = append(resp.Data, datastructs.SyncGet{
				Channel:  su.Channel,
				Name:     su.Name,
				Revision: su.Revision,
				Content:  string(su.Data),
			})
		}
		d, err := json.Marshal(resp)
		if err != nil {
			myutils.MyLogger.Println(err.Error())
			fmt.Fprintf(rw, ErrorResponse)
			return
		}
		myutils.MyLogger.Println("response: " + fmt.Sprintf("%+v", sus))
		rw.Write(d)
	} else {
		fmt.Fprintf(rw, ErrorResponse)
	}
}

// TODO: not battle proven yet
func Feedback(rw http.ResponseWriter, req *http.Request) {
	if !TheNexus.needFeedback {
		fmt.Fprintf(rw, ErrorResponse)
		return
	}

	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	myutils.MyLogger.Println("********File sync feedback********" + ip)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		myutils.MyLogger.Println(err.Error())
		fmt.Fprintf(rw, ErrorResponse)
		return
	}
	myutils.MyLogger.Println("request: " + string(body))

	var sfbs datastructs.SyncFeedbacks
	err = json.Unmarshal(body, &sfbs)
	if err != nil {
		myutils.MyLogger.Println(err.Error())
		fmt.Fprintf(rw, ErrorResponse)
		return
	}

	for _, sfb := range sfbs.Data {
		TheConsensus.candidates <- sfb
	}

	fmt.Fprintf(rw, OkResponse)
}

func Backup(rw http.ResponseWriter, req *http.Request) {
	backupRoot := req.FormValue("root")
	if backupRoot != "" {
		err := TheNexus.SaveChannels(backupRoot)
		if err != nil {
			myutils.MyLogger.Println(err.Error())
		}
	}
}
