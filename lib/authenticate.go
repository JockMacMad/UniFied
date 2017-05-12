package lib

import (
	"context"
	"fmt"
)

const loginBasePath = "/api/login"
const logoutBasePath = "/api/logoff"

// AuthenticateService is an interface for interfacing with the Authentication
// endpoints of the UniFi API
// See: https://developers.digitalocean.com/documentation/v2/#account
type AuthenticateService interface {
	Login(context.Context, string, string) (*Authentication, *Response, error)
	Logout(context.Context) (*Authentication, *Response, error)
}

// AlarmsServiceOp handles communication with the Alarm related methods of
// the UniFi API.
type AuthenticateServiceOp struct {
	client *UniFiClient
}

var _ AuthenticateService = &AuthenticateServiceOp{}

type authenticationRoot struct {
	Authentication *Authentication `json:"alarm"`
}

// Authentication represents a UniFi User Login
type Authentication struct {
	UserName      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Status        string `json:"status,omitempty"`
	StatusMessage string `json:"status_message,omitempty"`
}

// Get an alarm by ID.
func (s *AuthenticateServiceOp) Login(ctx context.Context, username string, password string) (*Authentication, *Response, error) {
	path := fmt.Sprintf("%s", loginBasePath)
	authRoot := new(authenticationRoot)
	authRoot.Authentication = new(Authentication)
	authRoot.Authentication.UserName = username
	authRoot.Authentication.Password = password

	req, err := s.client.NewRequest(ctx, "POST", path, authRoot.Authentication)
	if err != nil {
		return nil, nil, err
	}

	responseRoot := new(authenticationRoot)
	resp, err := s.client.Do(req, responseRoot)
	if err != nil {
		return nil, resp, err
	}

	return responseRoot.Authentication, resp, err
}

// Get an alarm by ID.
func (s *AuthenticateServiceOp) Logout(ctx context.Context) (*Authentication, *Response, error) {
	path := fmt.Sprintf("%s/%d", logoutBasePath)
	req, err := s.client.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(authenticationRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Authentication, resp, err
}

func (r Authentication) String() string {
	return Stringify(r)
}
