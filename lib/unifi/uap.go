package unifi

import (
	"context"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

// UAPService is an interface for interfacing with the UAP specific Device
// endpoints of the UniFi API
type UAPService interface {
	DisableAP(ctx context.Context, macAddress string, disabled bool) (*UniFiCmdResp, *Response, error)
	IsLocating(ctx context.Context, macAddress string) (bool, error)
	RenameAP(ctx context.Context, macAddress string, newName string) (*UniFiCmdResp, *Response, error)
	RestartAP(ctx context.Context, macAddress string) (*UniFiCmdResp, *Response, error)
	SetLocate(ctx context.Context, macAddress string, enabled bool) (*UniFiCmdResp, *Response, error)
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

// Disables a UniFi Acess Point which remains visible on the network & in the UniFi Controller.
func (uap *UAPServiceOp) DisableAP(
	ctx context.Context,
	macAddress string,
	disable bool) (*UniFiCmdResp, *Response, error) {

	uuid, err := uap.client.Devices.GetUUIDFromMac(ctx, macAddress)
	if err != nil {
		log.Error(err)
	}
	uapCmd := new(UAPCmdDisableAP)
	path := fmt.Sprintf("%s/%s", *uap.client.buildURL(restDeviceCmdBasePath), uuid)
	uapCmd.Disabled = disable

	return uap.client.sendCmd(ctx, "PUT", path, uapCmd)
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
func (uap *UAPServiceOp) SetLocate(
	ctx context.Context,
	macAddress string,
	enabled bool) (*UniFiCmdResp, *Response, error) {

	path := fmt.Sprintf("%s/%s", *uap.client.buildURL(devMgrCmdBasePath), macAddress)
	uapCmd := new(UniFiCmd)
	uapCmd.MacAddress = macAddress
	// If enabled is true the command is 'set-locate' to start flashing the APs LED
	if enabled {
		uapCmd.Cmd = "set-locate"
	} else {
		// If disabled i.e. false the command is 'unset-locate' to stop flashing the APs LED
		uapCmd.Cmd = "unset-locate"
	}

	return uap.client.sendCmd(ctx, "POST", path, uapCmd)
}

// Restarts i.e. reboots, the UniFi AP locating function to On or Off i.e. it is flashing it's LED in order
// someone can physically locate it visibly.
// macAddress is the MAC Address of the AP to configure
func (uap *UAPServiceOp) RestartAP(ctx context.Context, macAddress string) (*UniFiCmdResp, *Response, error) {
	path := fmt.Sprintf("%s/%s", *uap.client.buildURL(devMgrCmdBasePath), macAddress)
	uapCmd := new(UniFiCmd)
	uapCmd.MacAddress = macAddress
	uapCmd.Cmd = "restart"

	return uap.client.sendCmd(ctx, "POST", path, uapCmd)
}

// Renames the UniFi AP
// macAddress is the MAC Address of the AP to configure
func (uap *UAPServiceOp) RenameAP(ctx context.Context, macAddress string, newName string) (*UniFiCmdResp, *Response, error) {
	uuid, err := uap.client.Devices.GetUUIDFromMac(ctx, macAddress)
	if err != nil {
		log.Error(err)
	}
	uapCmd := new(UAPCmdRenameAP)
	path := fmt.Sprintf("%s/%s", *uap.client.buildURL(restDeviceCmdBasePath), uuid)
	uapCmd.Name = newName

	return uap.client.sendCmd(ctx, "POST", path, uapCmd)
}