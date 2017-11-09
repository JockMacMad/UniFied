package unifi

type UniFiCmd struct {
	Cmd        string `json:"cmd"`
	MacAddress string `json:"mac, omitempty"`
	UUID       string `json:"ma_id, omitempty"`
}

type UAPCmdDisableAP struct {
	Disabled bool `json:"disabled"`
}

type UAPCmdRenameAP struct {
	Name string `json:"name"`
}

type UniFiCmdResp struct {
	Data []interface{} `json:"data"`
	Meta struct {
		Status string `json:"rc"`
	} `json:"meta"`
}
