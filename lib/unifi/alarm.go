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

// AlarmsService is an interface for interfacing with the Alarm
// endpoints of the UniFi API
// See: https://developers.digitalocean.com/documentation/v2/#account
type AlarmsService interface {
	List(context.Context, *ListOptions) ([]Alarm, *Response, error)
	Get(context.Context, int) (*Alarm, *Response, error)
}

// AlarmsServiceOp handles communication with the Alarm related methods of
// the UniFi API.
type AlarmsServiceOp struct {
	client *UniFiClient
}

var _ AlarmsService = &AlarmsServiceOp{}

type alarmsRoot struct {
	Alarms []Alarm `json:"data"`
}

type alarmRoot struct {
	Alarm *Alarm `json:""`
}

// Alarm represents a UniFi Network Alarm
type Alarm struct {
	UUID           string `json:"_id"`
	Archived       bool   `json:"archived,omitempty"`
	DateTime       string `json:"datetime,omitempty"`
	Essid          string `json:"essid,omitempty"`
	HandledAdminId string `json:"handled_admin_id,omitempty"`
	HandledTime    string `json:"handled_time,omitempty"`
	Key            string `json:"key,omitempty"`
	MacAddress     string `json:"mac,omitempty"`
	Message        string `json:"msg,omitempty"`
	Occurs         int    `json:"occurs,omitempty"`
	SiteId         string `json:"site_id,omitempty"`
	SubSystem      string `json:"subsystem,omitempty"`
	//Time            *Timestamp  `json:"time,omitempty"`
	//EmailVerified   bool   `json:"email_verified,omitempty"`
	//Status          string `json:"status,omitempty"`
	//StatusMessage   string `json:"status_message,omitempty"`
}

// List all alarms
func (s *AlarmsServiceOp) List(ctx context.Context, opt *ListOptions) ([]Alarm, *Response, error) {
	//path := alarmsBasePath
	path := *s.buildURL()
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(alarmsRoot)
	resp, err := s.client.Do(req, root)

	if err != nil {
		return nil, resp, err
	}
	//if l := root.Links; l != nil {
	//	resp.Links = l
	//}

	if s.client.Options.DbUsage.DbUsageEnabled {

		root = AlarmsDB(s, root)
	}

	//log.Debug(root.Alarms)
	return root.Alarms, resp, err
}

//
func AlarmsDB(s *AlarmsServiceOp, root *alarmsRoot) *alarmsRoot {
	alarmsColExists := false
	var alarmsDB *db.Col = nil

	// Get the list of existing Columns in the DB
	for _, column := range s.client.Options.DbUsage.UnifiedDB.AllCols() {
		// If the Alarms Column already exists set the bool saying it exists and break out of the loop
		if column == "Alarms" {
			alarmsColExists = true
			break
		}
	}

	// After iterating the list of existing DB Columns Alarms does not exist so create it
	if alarmsColExists == false {
		if err := s.client.Options.DbUsage.UnifiedDB.Create("Alarms"); err != nil {
			panic(err)
		}

		alarmsDB = s.client.Options.DbUsage.UnifiedDB.Use("Alarms")

		if err := alarmsDB.Index([]string{"UUID"}); err != nil {
			panic(err)
		}

		// Now it exists so set the flag to true
		alarmsColExists = true
	} else {
		// Alarms Col already exists so just Use the existing one
		alarmsDB = s.client.Options.DbUsage.UnifiedDB.Use("Alarms")
		// It exists so set the flag to true
		alarmsColExists = true
	}

	if alarmsColExists {
		for i := 1; i < len(root.Alarms); i += 1 {
			v := root.Alarms[i]
			var query interface{}
			key := fmt.Sprintf(`[{"eq": "%s", "in": ["UUID"]}]`, v.UUID)
			json.Unmarshal(
				[]byte(key), &query)

			queryResult := make(map[int]struct{}) // query result (document IDs) goes into map keys

			if err := db.EvalQuery(query, alarmsDB, &queryResult); err != nil {
				panic(err)
			}

			if len(queryResult) == 0 {
				docID, err := alarmsDB.Insert(structs.Map(v))
				if err != nil {
					panic(err)
				}
				if log.GetLevel() == log.InfoLevel {
					log.WithFields(log.Fields{
						"docID":      docID,
						"Alarm UUID": v.UUID,
					}).Info(fmt.Sprintf("Alarm inserted.", docID, v.UUID))
				}
				if log.GetLevel() == log.DebugLevel {
					log.WithFields(log.Fields{
						"docID": docID,
						"Alarm": v,
					}).Debug(fmt.Sprintf("Alarm inserted."))
				}
			} else {
				// Query result are document IDs
				for id := range queryResult {
					// To get query result document, simply read it
					readBack, err := alarmsDB.Read(id)
					if err != nil {
						panic(err)
					}
					log.Info(fmt.Sprintf("Query returned document %v\n", readBack))
				}
			}
		}
	} else {
		panic(fmt.Sprintf("DB: Alarms column does not exist. Should not be possible."))
	}
	return root
}

// Get an alarm by its unique UUID.
func (s *AlarmsServiceOp) Get(ctx context.Context, uuid int) (*Alarm, *Response, error) {
	if uuid < 1 {
		return nil, nil, NewArgError("id", "cannot be less than 1")
	}

	if s.client.Options.DbUsage.DbUsageEnabled {

	}

	//path := fmt.Sprintf("%s/%d", alarmsBasePath, id)
	path := *s.buildURLWithId(uuid)
	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(alarmRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Alarm, resp, err
}

func (r Alarm) String() string {
	return Stringify(r)
}

func (r Alarm) Keys() []string {
	var keys []string

	s := structs.New(&r)
	fields := s.Fields()

	for _, f := range fields {
		fmt.Printf("field name: %+v\n", f.Name())
	}
	return keys
}

func (r Alarm) Values() []string {
	var values []string
	return values
}

func (s *AlarmsServiceOp) buildURL() *string {
	var buffer bytes.Buffer
	buffer.WriteString(s.client.BaseURL.String())
	buffer.WriteString(*s.client.SiteName)
	buffer.WriteString(alarmsBasePath)
	path := buffer.String()
	return &path
}

func (s *AlarmsServiceOp) buildURLWithId(id int) *string {
	var buffer bytes.Buffer
	buffer.WriteString(*s.buildURL())
	buffer.WriteString("/")
	buffer.WriteString(strconv.Itoa(id))
	path := buffer.String()
	return &path
}
