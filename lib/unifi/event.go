package unifi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/HouzuoGuo/tiedot/db"
	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"strconv"
)

//const eventsBasePath = "/api/s/default/list/event"
const eventsBasePath = "/list/event"

// AccountService is an interface for interfacing with the Account
// endpoints of the DigitalOcean API
// See: https://developers.digitalocean.com/documentation/v2/#account
type EventsService interface {
	List(context.Context, *ListOptions) ([]Event, *Response, error)
	Get(context.Context, int) (*Event, *Response, error)
}

// AccountServiceOp handles communication with the Account related methods of
// the DigitalOcean API.
type EventsServiceOp struct {
	client *UniFiClient
}

type eventsRoot struct {
	Events []Event `json:"data"`
	//Links  *Links  `json:"links"`
}

type eventRoot struct {
	Event *Event //`json:"event"`
}

var _ EventsService = &EventsServiceOp{}

// Account represents a DigitalOcean Account
type Event struct {
	UUID            string `json:"_id"`
	Key             string `json:"key"`
	Message         string `json:"msg"`
	SiteId          string `json:"site_id"`
	SubSystem       string `json:"subsystem"`
	AccessPoint     string `json:"ap,omitempty"`
	AccessPointName string `json:"ap_name,omitempty"`
	AccessPointFrom string `json:"ap_from,omitempty"`
	AccessPointTo   string `json:"ap_to,omitempty"`
	Admin           string `json:"admin,omitempty"`
	ByteCount       uint64 `json:"bytes,omitempty"`
	Channel         string `json:"channel,omitempty"`
	ChannelFrom     string `json:"channel_from,omitempty"`
	ChannelTo       string `json:"channel_to,omitempty"`
	DateTime        string `json:"datetime,omitempty"`
	Duration        int    `json:"duration,omitempty"`
	Essid           string `json:"essid,omitempty"`
	Gateway         string `json:"gateway,omitempty"`
	GatewayName     string `json:"gateway_name,omitempty"`
	Guest           string `json:"guest,omitempty"`
	Hostname        string `json:"hostname,omitempty"`
	IpAddress       string `json:"ip,omitempty"`
	IsAdmin         bool   `json:"is_admin,omitempty"`
	MacAddress      string `json:"mac,omitempty"`
	Minutes         string `json:"minutes,omitempty"`
	Network         string `json:"network,omitempty"`
	NumSta          int    `json:"num_sta,omitempty"`
	Radio           string `json:"radio,omitempty"`
	RadioFrom       string `json:"radio_from,omitempty"`
	RadioTo         string `json:"radio_to,omitempty"`
	Ssid            string `json:"ssid,omitempty"`
	Switch          string `json:"sw,omitempty"`
	SwitchName      string `json:"sw_name,omitempty"`
	VersionFrom     string `json:"version_from,omitempty"`
	VersionTo       string `json:"version_to,omitempty"`
	VouchersCreated string `json:"num,omitempty"`
	User            string `json:"user,omitempty"`
	//Time            *Timestamp  `json:"time,omitempty"`
}

// List all alarms
func (s *EventsServiceOp) List(ctx context.Context, opt *ListOptions) ([]Event, *Response, error) {
	//path := eventsBasePath
	path := *s.buildURL()
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(eventsRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}
	//	if l := root.Links; l != nil {
	//		resp.Links = l
	//	}

	if s.client.Options.DbUsage.DbUsageEnabled {

		root = EventsDB(s, root)
	}

	return root.Events, resp, err
}
func EventsDB(s *EventsServiceOp, root *eventsRoot) *eventsRoot {
	eventsColExists := false
	var eventsDB *db.Col = nil

	// Get the list of existing Columns in the DB
	for _, column := range s.client.Options.DbUsage.UnifiedDB.AllCols() {
		// If the Events Column already exists set the bool saying it exists and break out of the loop
		if column == "Events" {
			eventsColExists = true
			break
		}
	}

	// After iterating the list of existing DB Columns Events does not exist so create it
	if eventsColExists == false {
		if err := s.client.Options.DbUsage.UnifiedDB.Create("Events"); err != nil {
			panic(err)
		}

		eventsDB = s.client.Options.DbUsage.UnifiedDB.Use("Events")

		if err := eventsDB.Index([]string{"UUID"}); err != nil {
			panic(err)
		}

		// Now it exists so set the flag to true
		eventsColExists = true
	} else {
		// Events Col already exists so just Use the existing one
		eventsDB = s.client.Options.DbUsage.UnifiedDB.Use("Events")
		// It exists so set the flag to true
		eventsColExists = true
	}

	if eventsColExists {
		for i := 1; i < len(root.Events); i += 1 {
			v := root.Events[i]
			var query interface{}
			key := fmt.Sprintf(`[{"eq": "%s", "in": ["UUID"]}]`, v.UUID)
			json.Unmarshal(
				[]byte(key), &query)

			queryResult := make(map[int]struct{}) // query result (document IDs) goes into map keys

			if err := db.EvalQuery(query, eventsDB, &queryResult); err != nil {
				panic(err)
			}

			if len(queryResult) == 0 {
				docID, err := eventsDB.Insert(structs.Map(v))
				if err != nil {
					panic(err)
				}
				log.Info(fmt.Sprintf("Event inserted %s for UUID: ", docID, v.UUID))
			} else {
				// Query result are document IDs
				for id := range queryResult {
					// To get query result document, simply read it
					readBack, err := eventsDB.Read(id)
					if err != nil {
						panic(err)
					}
					//fmt.Println("Query returned document\n", readBack)
					log.Info(fmt.Sprintf("Query returned document %v\n", readBack))
				}
			}
		}
	}
	return root
}

// Get an alarm by ID.
func (s *EventsServiceOp) Get(ctx context.Context, id int) (*Event, *Response, error) {
	if id < 1 {
		return nil, nil, NewArgError("id", "cannot be less than 1")
	}

	//path := fmt.Sprintf("%s/%d", eventsBasePath, id)
	path := *s.buildURLWithId(id)
	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(eventRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Event, resp, err
}

func (s *EventsServiceOp) buildURL() *string {
	var buffer bytes.Buffer
	buffer.WriteString(s.client.BaseURL.String())
	buffer.WriteString(*s.client.SiteName)
	buffer.WriteString(eventsBasePath)
	path := buffer.String()
	return &path
}

func (s *EventsServiceOp) buildURLWithId(id int) *string {
	var buffer bytes.Buffer
	buffer.WriteString(*s.buildURL())
	buffer.WriteString("/")
	buffer.WriteString(strconv.Itoa(id))
	path := buffer.String()
	return &path
}
