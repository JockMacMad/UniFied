package unifi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/HouzuoGuo/tiedot/db"
	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
)

const usersBasePath = "/api/s/default/list/user"

// UsersService is an interface for interfacing with the user
// endpoints of the UniFi API
// See: https://developers.digitalocean.com/documentation/v2/#account
type UsersService interface {
	List(context.Context, *ListOptions) ([]User, *Response, error)
	Get(context.Context, int) (*User, *Response, error)
}

// UsersServiceOp handles communication with the User related methods of
// the UniFi API.
type UsersServiceOp struct {
	client *UniFiClient
}

var _ UsersService = &UsersServiceOp{}

type usersRoot struct {
	Users []User `json:"data"`
	//Links  *Links  `json:"links"`
}

type userRoot struct {
	//User *User `json:"data"`
	User *User `json:""`
}

// User represents a UniFi Network User
type User struct {
	UUID       string `json:"_id"`
	isGuest    bool   `json:"is_guest,omitempty"`
	isWired    bool   `json:"is_wired,omitempty"`
	OUI        string `json:"oui,omitempty"`
	MacAddress string `json:"mac,omitempty"`
	SiteId     string `json:"site_id,omitempty"`
}

// List all users
func (s *UsersServiceOp) List(ctx context.Context, opt *ListOptions) ([]User, *Response, error) {
	path := usersBasePath
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(usersRoot)
	resp, err := s.client.Do(req, root)

	if err != nil {
		return nil, resp, err
	}
	//if l := root.Links; l != nil {
	//	resp.Links = l
	//}
	if s.client.Options.DbUsage.DbUsageEnabled {
		root = UsersDB(s, root)
	}

	log.Debug(root.Users)
	return root.Users, resp, err
}

func UsersDB(s *UsersServiceOp, root *usersRoot) *usersRoot {
	usersColExists := false
	var usersDB *db.Col = nil

	// Get the list of existing Columns in the DB
	for _, column := range s.client.Options.DbUsage.UnifiedDB.AllCols() {
		// If the Users Column already exists set the bool saying it exists and break out of the loop
		if column == "Users" {
			usersColExists = true
			break
		}
	}

	// After iterating the list of existing DB Columns Users does not exist so create it
	if usersColExists == false {
		if err := s.client.Options.DbUsage.UnifiedDB.Create("Users"); err != nil {
			panic(err)
		}

		usersDB = s.client.Options.DbUsage.UnifiedDB.Use("Users")

		if err := usersDB.Index([]string{"UUID"}); err != nil {
			panic(err)
		}

		// Now it exists so set the flag to true
		usersColExists = true
	} else {
		// Users Col already exists so just Use the existing one
		usersDB = s.client.Options.DbUsage.UnifiedDB.Use("Users")
		// It exists so set the flag to true
		usersColExists = true
	}

	if usersColExists {
		for i := 1; i < len(root.Users); i += 1 {
			v := root.Users[i]
			var query interface{}
			key := fmt.Sprintf(`[{"eq": "%s", "in": ["UUID"]}]`, v.UUID)
			json.Unmarshal(
				[]byte(key), &query)

			queryResult := make(map[int]struct{}) // query result (document IDs) goes into map keys

			if err := db.EvalQuery(query, usersDB, &queryResult); err != nil {
				panic(err)
			}

			if len(queryResult) == 0 {
				docID, err := usersDB.Insert(structs.Map(v))
				if err != nil {
					panic(err)
				}
				log.Info(fmt.Sprintf("User inserted %s for UUID: ", docID, v.UUID))
			} else {
				// Query result are document IDs
				for id := range queryResult {
					// To get query result document, simply read it
					readBack, err := usersDB.Read(id)
					if err != nil {
						panic(err)
					}
					fmt.Println("Query returned document\n", readBack)
					//log.Info(fmt.Sprintf("Query returned document %v\n", readBack))
				}
			}
		}
	}
	return root
}

// Get an user by ID.
func (s *UsersServiceOp) Get(ctx context.Context, id int) (*User, *Response, error) {
	if id < 1 {
		return nil, nil, NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d", usersBasePath, id)
	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(userRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.User, resp, err
}

func (r User) String() string {
	return Stringify(r)
}
