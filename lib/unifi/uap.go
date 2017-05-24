package unifi

import (
	"bytes"
	"context"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http/httputil"
)

const devMgrCmdBasePath = "/cmd/devmgr"
const restDeviceCmdBasePath = "/rest/device"
const updDeviceCmdBasePath = "/upd/device"

// UAPService is an interface for interfacing with the UAP specific Device
// endpoints of the UniFi API
type UAPService interface {
	DisableAP(ctx context.Context, macAddress string, disabled bool) (*UAPCmdResp, *Response, error)
	IsLocating(ctx context.Context, macAddress string) (bool, error)
	RestartAP(ctx context.Context, macAddress string) (*UAPCmdResp, *Response, error)
	SetLocate(ctx context.Context, macAddress string, enabled bool) (*UAPCmdResp, *Response, error)
}

// UAPServiceOp handles communication with the Alarm related methods of
// the UniFi API.
type UAPServiceOp struct {
	client *UniFiClient
}

type UAP struct {
	UUID       string `json:"_id"`
	MacAddress string `json:"mac,omitempty"`
}

var _ UAPService = &UAPServiceOp{}

//
func (uap *UAPServiceOp) DisableAP(ctx context.Context, macAddress string, disable bool) (*UAPCmdResp, *Response, error) {
	uuid, err := uap.client.Devices.GetUUIDFromMac(ctx, macAddress)
	if err != nil {
		log.Error(err)
	}
	uapCmd := new(UAPCmdDisableAP)
	path := fmt.Sprintf("%s/%s", *uap.buildURL(restDeviceCmdBasePath), uuid)
	uapCmd.Disabled = disable

	return uap.sendCmd(ctx, "PUT", path, uapCmd)
}

// Checks to see the UniFi AP has Locating enabled i.e. it is flashing it's LED in order
// someone can physically locate it visibly.
// return true if the LED is flashing otherwise returns false
func (uap *UAPServiceOp) IsLocating(ctx context.Context, macAddress string) (bool, error) {
	device, _, err := uap.client.Devices.GetByMac(ctx, macAddress)
	if err != nil {
		log.Error(err)
	}
	return device.IsLocating, err
}

// Sets the UniFi AP locating function to On or Off i.e. it is flashing it's LED in order
// someone can physically locate it visibly.
// macAddress is the MAC Address of the AP to configure
// enabled is true to flash the APs LED and false is to disable the locating function.
func (uap *UAPServiceOp) SetLocate(ctx context.Context, macAddress string, enabled bool) (*UAPCmdResp, *Response, error) {
	path := fmt.Sprintf("%s/%s", *uap.buildURL(devMgrCmdBasePath), macAddress)
	uapCmd := new(UAPCmd)
	uapCmd.MacAddress = macAddress
	// If enabled is true the command is 'set-locate' to start flashing the APs LED
	if enabled {
		uapCmd.Cmd = "set-locate"
	} else {
		// If disabled i.e. false the command is 'unset-locate' to stop flashing the APs LED
		uapCmd.Cmd = "unset-locate"
	}

	return uap.sendCmd(ctx, "POST", path, uapCmd)
}

// Restarts i.e. reboots, the UniFi AP locating function to On or Off i.e. it is flashing it's LED in order
// someone can physically locate it visibly.
// macAddress is the MAC Address of the AP to configure
func (uap *UAPServiceOp) RestartAP(ctx context.Context, macAddress string) (*UAPCmdResp, *Response, error) {
	path := fmt.Sprintf("%s/%s", *uap.buildURL(devMgrCmdBasePath), macAddress)
	uapCmd := new(UAPCmd)
	uapCmd.MacAddress = macAddress
	uapCmd.Cmd = "restart"

	return uap.sendCmd(ctx, "POST", path, uapCmd)
}

func (uap *UAPServiceOp) sendCmd(ctx context.Context, method string, path string, uapCmd interface{}) (*UAPCmdResp, *Response, error) {
	// Create the HTTP Request
	req, err := uap.client.NewRequest(ctx, method, path, uapCmd)
	// Save a copy of this request for debugging.
	{
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Debug(err)
		}
		log.Debug(string(requestDump))
	}
	if err != nil {
		log.Error(err)
	}
	// Create the Response object to hold the results
	root := new(UAPCmdResp)
	// Make the HTTP Request to the UniFi Controller
	resp, err := uap.client.Do(req, root)
	if err != nil {
		log.Error(err)
		return nil, resp, err
	}
	return root, resp, err
}

func (uap *UAPServiceOp) buildURL(basePath string) *string {
	var buffer bytes.Buffer
	buffer.WriteString(uap.client.BaseURL.String())
	buffer.WriteString(*uap.client.SiteName)
	buffer.WriteString(basePath)
	path := buffer.String()
	return &path
}
