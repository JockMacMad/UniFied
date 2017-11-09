package unifi

import (
	"context"
	"fmt"
)

// ClientService is an interface for interfacing with client devices i.e. devices that connect
// to the UAP provided networks
type ClientService interface {
	BlockClient(ctx context.Context, macAddress string, blocked bool) (*UniFiCmdResp, *Response, error)
	AuthorizeGuest(
		ctx context.Context,
		macAddress string,
		minutes string,
		up string,
		down string,
		mbytes string,
		apMacAddress string) (*UniFiCmdResp, *Response, error)
	UnauthorizeGuest(
		ctx context.Context,
		macAddress string) (*UniFiCmdResp, *Response, error)
}

// ClientServiceOp handles communication with the Alarm related methods of
// the UniFi API.
type ClientServiceOp struct {
	client *UniFiClient
}

type AuthorizeGuestCmd struct {
	UniFiCmd
	Minutes      string `json:"minutes"`
	Up           string `json:"up, omitempty"`
	Down         string `json:"down, omitempty"`
	MBytes       string `json:"bytes, omitempty"`
	APMacAddress string `json:"ap_mac, omitempty"`
}

var _ ClientService = &ClientServiceOp{}

func (client *ClientServiceOp) AuthorizeGuest(
	ctx context.Context,
	macAddress string,
	minutes string,
	up string,
	down string,
	mbytes string,
	apMacAddress string) (*UniFiCmdResp, *Response, error) {

	clientCmd := new(AuthorizeGuestCmd)
	path := fmt.Sprintf("%s", *client.client.buildURL(cmdStaMgrCmdBasePath))
	// Mandatory params
	clientCmd.MacAddress = macAddress
	clientCmd.Minutes = minutes
	// Optional params
	if len(up) > 0 {
		clientCmd.Up = up
	}
	if len(down) > 0 {
		clientCmd.Down = down
	}
	if len(mbytes) > 0 {
		clientCmd.MBytes = mbytes
	}
	if len(apMacAddress) > 0 {
		clientCmd.APMacAddress = apMacAddress
	}

	return client.client.sendCmd(ctx, "POST", path, clientCmd)
}
func (client *ClientServiceOp) UnauthorizeGuest(
	ctx context.Context,
	macAddress string) (*UniFiCmdResp, *Response, error) {
	clientCmd := new(UniFiCmd)
	path := fmt.Sprintf("%s", *client.client.buildURL(cmdStaMgrCmdBasePath))

	clientCmd.Cmd = "unauthorize-guest"
	clientCmd.MacAddress = macAddress

	return client.client.sendCmd(ctx, "POST", path, clientCmd)
}

func (client *ClientServiceOp) BlockClient(
	ctx context.Context,
	macAddress string,
	blocked bool) (*UniFiCmdResp, *Response, error) {

	clientCmd := new(UniFiCmd)
	path := fmt.Sprintf("%s", *client.client.buildURL(cmdStaMgrCmdBasePath))
	if blocked {
		clientCmd.Cmd = "block-sta"
	} else {
		clientCmd.Cmd = "unblock-sta"
	}
	clientCmd.MacAddress = macAddress

	return client.client.sendCmd(ctx, "POST", path, clientCmd)
}
