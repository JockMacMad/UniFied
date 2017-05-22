package unifi

type UAPCmd struct {
	Cmd		string		`json:"cmd"`
	MacAddress	string		`json:"mac"`
}

type UAPCmdResp struct {
	Data []interface{} `json:"data"`
	Meta struct {
		Status string `json:"rc"`
	} `json:"meta"`
}