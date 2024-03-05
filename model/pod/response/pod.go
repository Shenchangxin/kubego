package response

//@Author: morris
//NAME 名称
//READY   状态 1/2
//STATUS   Running/Error
//RESTARTS  n次
//AGE  运行时间
//IP   podid
//NODE pod被调度到哪个node
type PodListItem struct {
	Name     string `json:"name"`
	Ready    string `json:"ready"`
	Status   string `json:"status"`
	Restarts int32  `json:"restarts"`
	Age      int64  `json:"age"`
	IP       string `json:"IP"`
	Node     string `json:"node"`
}
