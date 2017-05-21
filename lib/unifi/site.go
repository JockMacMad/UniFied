package unifi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/HouzuoGuo/tiedot/db"
	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
)

const sitesBasePath = "/v1/sites"

// SitesService is an interface for interfacing with the Site
// endpoints of the UniFi API
// See: https://developers.digitalocean.com/documentation/v2/#account
type SitesService interface {
	List(context.Context, *ListOptions) ([]Site, *Response, error)
	Get(context.Context, int) (*Site, *Response, error)
}

// SitesServiceOp handles communication with the Site related methods of
// the UniFi API.
type SitesServiceOp struct {
	client *UniFiClient
}

var _ SitesService = &SitesServiceOp{}

type sitesRoot struct {
	Sites []Site `json:"sites"`
}

type siteRoot struct {
	Site *Site `json:"site"`
}

// Site represents a UniFi Network Site
type Site struct {
	FloatingIPLimit int    `json:"floating_ip_limit,omitempty"`
	Email           string `json:"email,omitempty"`
	UUID            string `json:"uuid,omitempty"`
	EmailVerified   bool   `json:"email_verified,omitempty"`
	Status          string `json:"status,omitempty"`
	StatusMessage   string `json:"status_message,omitempty"`
}

// List all Sites
func (s *SitesServiceOp) List(ctx context.Context, opt *ListOptions) ([]Site, *Response, error) {
	path := sitesBasePath
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(sitesRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}
	//	if l := root.Links; l != nil {
	//		resp.Links = l
	//	}

	if s.client.Options.DbUsage.DbUsageEnabled {
		root = SitesDB(s, root)
	}

	return root.Sites, resp, err
}

func SitesDB(s *SitesServiceOp, root *sitesRoot) *sitesRoot {

	if err := s.client.Options.DbUsage.UnifiedDB.Create("Sites"); err != nil {
		panic(err)
	}

	sitesDB := s.client.Options.DbUsage.UnifiedDB.Use("Sites")
	if err := sitesDB.Index([]string{"UUID"}); err != nil {
		panic(err)
	}

	for i := 1; i < len(root.Sites); i += 1 {
		v := root.Sites[i]
		var query interface{}
		key := fmt.Sprintf(`[{"eq": "%s", "in": ["UUID"]}]`, v.UUID)
		json.Unmarshal(
			[]byte(key), &query)

		queryResult := make(map[int]struct{}) // query result (document IDs) goes into map keys

		if err := db.EvalQuery(query, sitesDB, &queryResult); err != nil {
			panic(err)
		}

		docID, err := sitesDB.Insert(structs.Map(v))
		if err != nil {
			panic(err)
		}
		log.Debugln(docID)
	}
	return root
}

// Get an Site by ID.
func (s *SitesServiceOp) Get(ctx context.Context, id int) (*Site, *Response, error) {
	if id < 1 {
		return nil, nil, NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d", sitesBasePath, id)
	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(siteRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Site, resp, err
}

func (r Site) String() string {
	return Stringify(r)
}
