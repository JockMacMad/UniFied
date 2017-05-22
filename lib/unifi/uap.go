package unifi

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"bytes"
	"fmt"
	"net/http/httputil"
)

const devMgrCmdBasePath = "/cmd/devmgr"

// UAPService is an interface for interfacing with the UAP specific Device
// endpoints of the UniFi API
type UAPService interface {
	SetLocate(ctx context.Context, macAddress string, enabled bool) (*UAPCmdResp, *Response, error)
	IsLocating(ctx context.Context, macAddress string) (bool, error)
}

// UAPServiceOp handles communication with the Alarm related methods of
// the UniFi API.
type UAPServiceOp struct {
	client 	*UniFiClient
}

type UAP struct {
	MacAddress             string          `json:"mac,omitempty"`
}

var _ UAPService = &UAPServiceOp{}

// Sets the UniFi AP locating function to On or Off i.e. it is flashing it's LED in order
// someone can physically locate it visibly.
// macAddress is the MAC Address of the AP to configure
// enabled is true to flash the APs LED and false is to disable the locating function.
func (uap *UAPServiceOp) SetLocate(ctx context.Context, macAddress string, enabled bool) (*UAPCmdResp, *Response, error) {
	path := fmt.Sprintf("%s/%s", *uap.buildURL(), macAddress)
	uapCmd := new(UAPCmd)
	uapCmd.MacAddress = macAddress
	// If enabled is true the command is 'set-locate' to start flashing the APs LED
	if enabled {
		uapCmd.Cmd = "set-locate"
	} else {
		// If disabled i.e. false the command is 'unset-locate' to stop flashing the APs LED
		uapCmd.Cmd = "unset-locate"
	}

	// Create the HTTP Request
	req, err := uap.client.NewRequest(ctx, "POST", path, uapCmd)
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

func (uap *UAPServiceOp) buildURL() *string {
	var buffer bytes.Buffer
	buffer.WriteString(uap.client.BaseURL.String())
	buffer.WriteString(*uap.client.SiteName)
	buffer.WriteString(devMgrCmdBasePath)
	path := buffer.String()
	return &path
}