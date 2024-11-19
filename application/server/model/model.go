package model

// 定义用户结构体
type User struct {
	UserRSAPK string   `json:"userRSAPK"`
	UserName  string   `json:"userName"`
	Videos    []*Video `json:"video"`
}

// 定义作品结构体
type Video struct {
	Owner     string `json:"owner"`
	OrgInfo   string `json:"orgInfo"`
	VideoCID  string `json:"videoCID"`
	VideoName string `json:"videoName"`
	Txid      string `json:"txid"`
	Timestamp string `json:"timestamp"`
}
