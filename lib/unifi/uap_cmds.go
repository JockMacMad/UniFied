package unifi

type UAPCmd struct {
	Cmd        string `json:"cmd"`
	MacAddress string `json:"mac, omitempty"`
	UUID       string `json:"ma_id, omitempty"`
}

type UAPCmdDisableAP struct {
	Disabled bool `json:"disabled"`
}

type UAPCmdResp struct {
	Data []interface{} `json:"data"`
	Meta struct {
		Status string `json:"rc"`
	} `json:"meta"`
}
