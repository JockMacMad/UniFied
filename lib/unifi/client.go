package unifi

import (
	"context"
	"fmt"
)

// ClientService is an interface for interfacing with client devices i.e. devices that connect
// to the UAP provided networks
type ClientService interface {
	BlockClient(ctx context.Context, macAddress string, blocked bool) (*UAPCmdResp, *Response, error)
}

// ClientServiceOp handles communication with the Alarm related methods of
// the UniFi API.
type ClientServiceOp struct {
	client *UniFiClient
}

type Client struct {
	MacAddress string `json:"mac,omitempty"`
	IsBlocked  bool   `json:"mac,omitempty"`
}

var _ ClientService = &ClientServiceOp{}

func (client *ClientServiceOp) BlockClient(
	ctx context.Context,
	macAddress string,
	blocked bool) (*UAPCmdResp, *Response, error) {

	clientCmd := new(Client)
	path := fmt.Sprintf("%s", *client.client.buildURL(restDeviceCmdBasePath))
	clientCmd.IsBlocked = blocked

	return client.client.sendCmd(ctx, "POST", path, clientCmd)
}
