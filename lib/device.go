package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/HouzuoGuo/tiedot/db"
	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
)

const devicesBasePath = "/api/s/default/stat/device"

// DeviceService is an interface for interfacing with the Device
// endpoints of the UniFi API
// See: https://developers.digitalocean.com/documentation/v2/#account
type DevicesService interface {
	List(context.Context, *ListOptions) ([]Device, *Response, error)
	ListShort(context.Context, string, *ListOptions) ([]DeviceShort, *Response, error)
	Get(context.Context, int) (*Device, *Response, error)
	GetByMac(context.Context, string) (*Device, *Response, error)
	GetIPFromMac(ctx context.Context, mac string) (string, error)
}

// AlarmsServiceOp handles communication with the Alarm related methods of
// the UniFi API.
type DevicesServiceOp struct {
	client *UniFiClient
}

type devicesRoot struct {
	Devices []Device `json:"data"`
}

type deviceRoot struct {
	Device []Device `json:"data"`
}

var _ DevicesService = &DevicesServiceOp{}

// Device represents a UniFi Network Device
type Device struct {
	UUID                   string          `json:"_id"`
	IsAdopted              bool            `json:"adopted,omitempty"`
	BoardRev               int             `json:"board_rev,omitempty"`
	Bytes                  int             `json:"bytes,omitempty"`
	ConfigVersion          string          `json:"cfgversion,omitempty"`
	ConfigNetwork          ConfigNetwork   `json:"config_network,omitempty"`
	ConnectRequestIP       string          `json:"connect_request_ip,omitempty"`
	ConnectRequestPort     string          `json:"connect_request_port,omitempty"`
	DeviceId               string          `json:"device_id,omitempty"`
	DhcpServerTable        []string        `json:"dhcp_server_table,omitempty"`
	IsDot1xPortCtrlEnabled bool            `json:"dot1x_portctrl_enabled,omitempty"`
	DownLinks              []DownLinkTable `json:"downlink_table,omitempty"`
	Ethernet               []EthernetTable `json:"ethernet_table,omitempty"`
	IsFlowControlEnabled   bool            `json:"flowctrl_enabled,omitempty"`
	FWCaps                 int             `json:"fw_caps,omitempty"`
	GeneralTemperature     int             `json:"general_temperature,omitempty"`
	GuestNumSta            int             `json:"guest-num_sta,omitempty"`
	HasFan                 bool            `json:"has_fan,omitempty"`
	InformAuthkey          string          `json:"inform_authkey,omitempty"`
	InformIP               string          `json:"inform_ip,omitempty"`
	InformURL              string          `json:"inform_url,omitempty"`
	IP                     string          `json:"ip,omitempty"`
	AreJumboFrameEnabled   bool            `json:"jumboframe_enabled,omitempty"`
	MacAddress             string          `json:"mac,omitempty"`
	Name                   string          `json:"name,omitempty"`
	Model                  string          `json:"model,omitempty"`
	KnownCfgVersion        string          `json:"known_cfgversion,omitempty"`
	LEDOverride            string          `json:"led_override,omitempty"`
	IsLocating             bool            `json:"locating,omitempty"`
	IsOverHeating          bool            `json:"overheating,omitempty"`
	PortOverrides          []PortOverrides `json:"port_overrides,omitempty"`
	Ports                  []PortTable     `json:"port_table,omitempty"`
	Serial                 string          `json:"serial,omitempty"`
	SiteId                 string          `json:"site_id,omitempty"`
	State                  int             `json:"state,omitempty"`
	STPPriority            string          `json:"stp_priority,omitempty"`
	STPVersion             string          `json:"stp_version,omitempty"`
	Type                   string          `json:"type,omitempty"`
	UplinkDepth            int             `json:"uplink_depth,omitempty"`
	Version                string          `json:"version,omitempty"`
	//Time            *Timestamp  `json:"time,omitempty"`
}

type DeviceShort struct {
	Type       string `json:"type,omitempty"`
	Serial     string `json:"serial,omitempty"`
	SiteId     string `json:"site_id,omitempty"`
	MacAddress string `json:"mac,omitempty"`
	Name       string `json:"name,omitempty"`
	Model      string `json:"model,omitempty"`
	IP         string `json:"ip,omitempty"`
	UUID       string `json:"_id"`
	IsAdopted  bool   `json:"adopted,omitempty"`
	Version    string `json:"version,omitempty"`
}

type ConfigNetwork struct {
	Dns1      string `json:"dns1,omitempty"`
	Dns2      string `json:"dns2,omitempty"`
	DnsSuffix string `json:"dnssuffix,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	IP        string `json:"ip,omitempty"`
	Netmask   string `json:"netmask,omitempty"`
	Type      string `json:"type,omitempty"`
}

type DownLinkTable struct {
	IsFullDuplex bool   `json:"full_duplex,omitempty"`
	MacAddress   string `json:"mac"`
	PortIdx      int    `json:"port_idx,omitempty"`
	Speed        int    `json:"speed,omitempty"`
}

type EthernetTable struct {
	MacAddress string `json:"mac"`
	Name       string `json:"name,omitempty"`
	NumPort    int    `json:"num_port,omitempty"`
}

type PortOverrides struct {
	IsAutoNeg               bool   `json:"autoneg,omitempty"`
	IsIsolation             bool   `json:"isolation,omitempty"`
	Name                    string `json:"name,omitempty"`
	OpMode                  string `json:"op_mode,omitempty"`
	POEMode                 string `json:"poe_mode,omitempty"`
	PortIdx                 int    `json:"port_idx,omitempty"`
	PortConfId              string `json:"portconf_id,omitempty"`
	IsStormCtrlBcastEnabled bool   `json:"stormctrl_bcast_enabled,omitempty"`
	IsStormCtrlMcastEnabled bool   `json:"stormctrl_mcast_enabled,omitempty"`
	IsStormCtrlUcastEnabled bool   `json:"stormctrl_ucast_enabled,omitempty"`
}

type PortTable struct {
	IsAggregatedBy          bool   `json:"aggregated_by,omitempty"`
	IsAutoNegEnabled        bool   `json:"autoneg,omitempty"`
	BytesR                  int64  `json:"bytes-r,omitempty"`
	Dot1xMode               string `json:"dot1x_mode,omitempty"`
	Dot1xStatus             string `json:"dot1x_status,omitempty"`
	IsEnabled               bool   `json:"enable,omitempty"`
	IsFlowControlRXEnabled  bool   `json:"flowctrl_rx,omitempty"`
	IsFlowControlTXEnabled  bool   `json:"flowctrl_tx,omitempty"`
	IsFullDuplexEnabled     bool   `json:"full_duplex,omitempty"`
	IsUplink                bool   `json:"is_uplink,omitempty"`
	IsIsolated              bool   `json:"isolation,omitempty"`
	IsJumboEnabled          bool   `json:"jumbo,omitempty"`
	IsMasked                bool   `json:"masked,omitempty"`
	Media                   string `json:"media,omitempty"`
	Name                    string `json:"name,omitempty"`
	OpMode                  string `json:"op_mode,omitempty"`
	POECaps                 int    `json:"poe_caps,omitempty"`
	POEClass                string `json:"poe_class,omitempty"`
	POECurrent              string `json:"poe_current,omitempty"`
	IsPOEEnabled            bool   `json:"poe_enable,omitempty"`
	IsPOEGood               bool   `json:"poe_good,omitempty"`
	POEMode                 string `json:"poe_mode,omitempty"`
	POEPower                string `json:"poe_power,omitempty"`
	POEVoltage              string `json:"poe_voltage,omitempty"`
	PortIdx                 int    `json:"port_idx,omitempty"`
	IsPortPOE               bool   `json:"port_poe,omitempty"`
	PortConfId              string `json:"portconf_id,omitempty"`
	RXBroadcast             int64  `json:"rx_broadcast,omitempty"`
	RXBytes                 int64  `json:"rx_bytes,omitempty"`
	RXBytesR                int64  `json:"rx_bytes-r,omitempty"`
	RXDropped               int64  `json:"rx_dropped,omitempty"`
	RXErrors                int64  `json:"rx_errors,omitempty"`
	RXMulticast             int64  `json:"rx_multicast,omitempty"`
	RXPackets               int64  `json:"rx_packets,omitempty"`
	PortSpeed               int    `json:"speed,omitempty"`
	IsStormCtrlBcastEnabled bool   `json:"stormctrl_bcast_enabled,omitempty"`
	IsStormCtrlMcastEnabled bool   `json:"stormctrl_mcast_enabled,omitempty"`
	IsStormCtrlUcastEnabled bool   `json:"stormctrl_ucast_enabled,omitempty"`
	STPPathCost             int    `json:"stp_pathcost,omitempty"`
	STPState                string `json:"stp_state,omitempty"`
	TXBroadcast             int64  `json:"tx_broadcast,omitempty"`
	TXBytes                 int64  `json:"tx_bytes,omitempty"`
	TXBytesR                int64  `json:"tx_bytes-r,omitempty"`
	TXDropped               int64  `json:"tx_dropped,omitempty"`
	TXErrors                int64  `json:"tx_errors,omitempty"`
	TXMulticast             int64  `json:"tx_multicast,omitempty"`
	TXPackets               int64  `json:"tx_packets,omitempty"`
	IsUp                    bool   `json:"up,omitempty"`
}

type PortStats struct {
}

type SysStats struct {
}

// List all devices
func (s *DevicesServiceOp) List(ctx context.Context, opt *ListOptions) ([]Device, *Response, error) {
	path := devicesBasePath
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(devicesRoot)
	resp, err := s.client.Do(req, root)

	if err != nil {
		return nil, resp, err
	}
	//if l := root.Links; l != nil {
	//	resp.Links = l
	//}

	if s.client.Options.DbUsage.DbUsageEnabled {

		root = DevicesDB(s, root)
	}

	//log.Debug(root.Devices)
	return root.Devices, resp, err
}

// List all devices
func (s *DevicesServiceOp) ListShort(ctx context.Context, filter string, opt *ListOptions) ([]DeviceShort, *Response, error) {
	path := devicesBasePath
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(devicesRoot)
	resp, err := s.client.Do(req, root)

	if err != nil {
		return nil, resp, err
	}

	if s.client.Options.DbUsage.DbUsageEnabled {

		root = DevicesDB(s, root)
	}

	var deviceShortArray []DeviceShort
	for _, dev := range root.Devices {
		switch filter {
		case dev.Type:
			deviceShortArray = append(deviceShortArray, dev.toDeviceShort())
		case "all":
			deviceShortArray = append(deviceShortArray, dev.toDeviceShort())
		}
	}

	//log.Debug(root.Devices)
	return deviceShortArray, resp, err
}
func DevicesDB(s *DevicesServiceOp, root *devicesRoot) *devicesRoot {
	devicesColExists := false
	var devicesDB *db.Col = nil

	// Get the list of existing Columns in the DB
	for _, column := range s.client.Options.DbUsage.UnifiedDB.AllCols() {
		// If the Devices Column already exists set the bool saying it exists and break out of the loop
		if column == "Devices" {
			devicesColExists = true
			break
		}
	}

	// After iterating the list of existing DB Columns Devices does not exist so create it
	if devicesColExists == false {
		if err := s.client.Options.DbUsage.UnifiedDB.Create("Devices"); err != nil {
			panic(err)
		}

		devicesDB = s.client.Options.DbUsage.UnifiedDB.Use("Devices")

		if err := devicesDB.Index([]string{"UUID"}); err != nil {
			panic(err)
		}

		// Now it exists so set the flag to true
		devicesColExists = true
	} else {
		// Devices Col already exists so just Use the existing one
		devicesDB = s.client.Options.DbUsage.UnifiedDB.Use("Devices")
		// It exists so set the flag to true
		devicesColExists = true
	}

	if devicesColExists {
		for i := 1; i < len(root.Devices); i += 1 {
			v := root.Devices[i]
			var query interface{}
			key := fmt.Sprintf(`[{"eq": "%s", "in": ["UUID"]}]`, v.UUID)
			json.Unmarshal(
				[]byte(key), &query)

			queryResult := make(map[int]struct{}) // query result (document IDs) goes into map keys

			if err := db.EvalQuery(query, devicesDB, &queryResult); err != nil {
				panic(err)
			}

			if len(queryResult) == 0 {
				docID, err := devicesDB.Insert(structs.Map(v))
				if err != nil {
					panic(err)
				}
				if log.GetLevel() == log.InfoLevel {
					log.WithFields(log.Fields{
						"docID":       docID,
						"Device UUID": v.UUID,
					}).Info(fmt.Sprintf("Device inserted.", docID, v.UUID))
				}
				if log.GetLevel() == log.DebugLevel {
					log.WithFields(log.Fields{
						"docID":  docID,
						"Device": v,
					}).Debug(fmt.Sprintf("Device inserted."))
				}
			} else {
				// Query result are document IDs
				for id := range queryResult {
					// To get query result document, simply read it
					readBack, err := devicesDB.Read(id)
					if err != nil {
						panic(err)
					}
					log.Info(fmt.Sprintf("Query returned document %v\n", readBack))
				}
			}
		}
	} else {
		panic(fmt.Sprintf("DB: Devices column does not exist. Should not be possible."))
	}
	return root
}

// Get an Device by ID.
func (s *DevicesServiceOp) Get(ctx context.Context, id int) (*Device, *Response, error) {
	if id < 1 {
		return nil, nil, NewArgError("id", "cannot be less than 1")
	}

	if s.client.Options.DbUsage.DbUsageEnabled {

	}

	path := fmt.Sprintf("%s/%d", devicesBasePath, id)
	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(devicesRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return &root.Devices[0], resp, err
}

// Get an Device by ID.
func (s *DevicesServiceOp) GetByMac(ctx context.Context, mac string) (*Device, *Response, error) {

	if s.client.Options.DbUsage.DbUsageEnabled {

	}

	path := fmt.Sprintf("%s/%s", devicesBasePath, mac)
	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(devicesRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return &root.Devices[0], resp, err
}

func (d Device) String() string {
	return Stringify(d)
}

func (d Device) Keys() []string {
	var keys []string

	s := structs.New(&d)
	fields := s.Fields()

	for _, f := range fields {
		fmt.Printf("field name: %+v\n", f.Name())
	}
	return keys
}

func (d Device) Values() []string {
	var values []string
	return values
}

func (d Device) toDeviceShort() DeviceShort {
	shortStruct := DeviceShort{Name: d.Name, UUID: d.UUID, Version: d.Version, SiteId: d.SiteId,
		MacAddress: d.MacAddress, IP: d.IP, IsAdopted: d.IsAdopted, Model: d.Model, Serial: d.Serial,
		Type: d.Type,
	}
	return shortStruct
}

func (s *DevicesServiceOp) GetIPFromMac(ctx context.Context, mac string) (string, error) {
	device, _, err := s.GetByMac(ctx, mac)
	if err != nil {
		return "", err
	}
	return device.IP, nil
}

//func GetDeviceFromMac(mac string, short bool) Device {

//}
