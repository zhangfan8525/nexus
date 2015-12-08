package datastructs

type Status struct {
	Status int `json:"status"`
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
