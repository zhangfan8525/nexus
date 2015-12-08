package datastructs

type Update struct {
	Id          int64 `json:"id"`
	Update_type int   `json:"type"`
	Start_coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"start_coord"`
	End_coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"end_coord"`
	Node_coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"node_coord,omitempty"`
	Via        [][]float64 `json:"via"`
	Start_time string      `json:"start_time,omitempty"`
	End_time   string      `json:"end_time,omitempty"`
	Repeat     [][2]int    `json:"repeat,omitempty"`
	Add_time   string      `json:"add_time,omitempty"`
	Value      string      `json:"value"`
	Status     int         `json:"status"`
	Distance   float64     `json:"distance,omitempty"`
}

type Patch struct {
	Id          int64  `json:"id"`
	Update_time string `json:"update_time,omitempty"`
	Version     int    `json:"version,omitempty"`
	Status      int    `json:"status"`
	Patch       string `json:"patch,omitempty"`
	Add_time    string `json:"add_time,omitempty"`
	Error_count int    `json:"error_count,omitempty"`
}

type Patches struct {
	Status      int      `json:"status"`
	Data        []Update `json:"patch"`
	Update_time string   `json:"update_time"`
}

type Query struct {
	Patch_id  []int64 `json:"patch_id"`
	Type      int     `json:"type"`
	From_time string  `json:"from_time"`
	To_time   string  `json:"to_time"`
	Top_right struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"top_right"`
	Lower_left struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"lower_left"`
}

type Status struct {
	Status int `json:"status"`
}

type InsertResp struct {
	Status int   `json:"status"`
	Id     int64 `json:"id"`
}

type Monitor struct {
	Desc   string `json:"Desc"`
	Ip     string `json:"IP"`
	Id     string `json:"Id"`
	Status string `json:"Status"`
	Time   string `json:"Time"`
}

type Inspect struct {
	Status    int     `json:"status"`
	Count     int     `json:"count"`
	Version   string  `json:"version"`
	Boot_time string  `json:"boot_time"`
	Id_list   []int64 `json:"id_list"`
}

type SyncPut struct {
	Name    string `json:"name"`
	Channel string `json:"channel"`
	Content string `json:"content"`
}

type SyncFeedback struct {
	Channel  string `json:"channel"`
	Name     string `json:"name"`
	Revision int    `json:"revision"`
	Valid    bool   `json:"valid"`
}

type SyncFeedbacks struct {
	Data []SyncFeedback `json:"data"`
}

type SyncGet struct {
	Channel  string `json:"channel"`
	Name     string `json:"name"`
	Revision int    `json:"revision"`
	Content  string `json:"content"`
}

type SyncGetResp struct {
	Status int       `json:"status"`
	Data   []SyncGet `json:"data"`
}

type SyncGetReq struct {
	Channel   string `json:"channel"`
	TimeStamp int64  `json:"timestamp"`
}
