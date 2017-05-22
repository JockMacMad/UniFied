package unifi

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/HouzuoGuo/tiedot/db"
	log "github.com/Sirupsen/logrus"
	"github.com/google/go-querystring/query"
	"github.com/logmatic/logmatic-go"
	"github.com/shiena/ansicolor"
	headerLink "github.com/tent/http-link-go"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"strconv"
)

const (
	libraryVersion = "0.0.1"
	defaultBaseURL = "https://192.168.10.7:8443"
	userAgent      = "unified/" + libraryVersion
	mediaType      = "application/json"

	headerRateLimit     = "RateLimit-Limit"
	headerRateRemaining = "RateLimit-Remaining"
	headerRateReset     = "RateLimit-Reset"
)

// Simple structure that captures the URL (string) and Port (int) that the DB is running on.
type UnifiedDBHost struct {
	dbUrl  string
	dbPort int
}

// Simple structure used to hold information on the usage of a DB with Unified to store information retreived from
// the Ubiquiti UniFi Controller.
type UnifiedDBOptions struct {
	// Actual pointer to the DB
	UnifiedDB *db.DB
	// Should we use a DB or be transient to the UniFi data.
	DbUsageEnabled bool
	// If we are using a DB are we to use an in memory DB?
	UseInMemoryDB bool
	// If we are using a DB but not an in-memory DB what is the
	//   DB Host information.
	dbHost *UnifiedDBHost
}

type UnifiedOptions struct {
	DbUsage *UnifiedDBOptions
}

// The main UniFi Client structure. This holds pointers to the actual HTTP Client, Unifi Controller Username & Password,
// Site etc. to use.
type UniFiClient struct {
	// HTTP client used to communicate with the DO API.
	client *http.Client

	// User set options on the UniFiClient
	Options *UnifiedOptions

	// The user using Unified
	UserName *string

	// The password for the user using Unified
	Password *string

	UnifiCookie *http.Cookie
	CSRFCookie  *http.Cookie

	// Specified site to operate on
	SiteName *string

	// Base URL for API requests.
	BaseURL *url.URL

	// SSH Client
	SSHClient  *ssh.Client
	SSHSession *ssh.Session

	// User agent for client
	UserAgent string

	// Rate contains the current rate limit for the client as determined by the most recent
	// API call.
	Rate Rate

	// Services used for communicating with the API
	Alarms         AlarmsService
	Authentication AuthenticateService
	Devices        DevicesService
	Events         EventsService
	Sites          SitesService
	Users          UsersService
	UAP		UAPService

	// Optional function called after every successful request made to the DO APIs
	onRequestCompleted RequestCompletionCallback
}

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response)

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`
}

// Response is a DigitalOcean response. This wraps the standard http.Response returned from DigitalOcean.
type Response struct {
	*http.Response
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message
	Message string `json:"message"`

	// RequestID returned from the API, useful to contact support.
	RequestID string `json:"request_id"`
}

// Rate contains the rate limit for the current client.
type Rate struct {
	// The number of request per hour the client is currently limited to.
	Limit int `json:"limit"`

	// The number of remaining requests the client can make this hour.
	Remaining int `json:"remaining"`

	// The time at which the current rate limit will reset.
	//Reset Timestamp `json:"reset"`
}

// SetBaseURL is a client option for setting the base URL.
func SetBaseURL(bu string) ClientOpt {
	return func(c *UniFiClient) error {
		u, err := url.Parse(bu)
		if err != nil {
			return err
		}

		c.BaseURL = u
		return nil
	}
}

// SetUserAgent is a client option for setting the user agent.
func SetUserAgent(ua string) ClientOpt {
	return func(c *UniFiClient) error {
		c.UserAgent = fmt.Sprintf("%s+%s", ua, c.UserAgent)
		return nil
	}
}

// NewUniFiClient returns a new DigitalOcean API client.
func NewUniFiClient(httpClient *http.Client, options *UnifiedOptions) *UniFiClient {
	f, err := os.OpenFile("unified.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {

	}
	log.SetOutput(f)
	log.SetFormatter(&logmatic.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.Println("Unified Daemon Starting... ")
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &UniFiClient{client: httpClient, Options: options, BaseURL: baseURL, UserAgent: userAgent}
	c.Alarms = &AlarmsServiceOp{client: c}
	c.Authentication = &AuthenticateServiceOp{client: c}
	c.Devices = &DevicesServiceOp{client: c}
	c.Events = &EventsServiceOp{client: c}
	c.Users = &UsersServiceOp{client: c}
	c.UAP = &UAPServiceOp{client: c}

	if c.Options.DbUsage.DbUsageEnabled {
		unfiedDBLocation := "/tmp/UnifiedDB"
		os.RemoveAll(unfiedDBLocation)
		defer os.RemoveAll(unfiedDBLocation)

		// (Create if not exist) open a database
		unifiedDB, err := db.OpenDB(unfiedDBLocation)
		if err != nil {
			panic(err)
		}
		c.Options.DbUsage.UnifiedDB = unifiedDB
	}

	return c
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()
	return origURL.String(), nil
}

// ClientOpt are options for New.
type ClientOpt func(*UniFiClient) error

// New returns a new Unified API client instance.
func New(httpClient *http.Client, o *UnifiedOptions, opts ...ClientOpt) (*UniFiClient, error) {
	c := NewUniFiClient(httpClient, o)
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the UniFi Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *UniFiClient) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	if c.UnifiCookie != nil {
		req.AddCookie(c.UnifiCookie)
	}
	if c.CSRFCookie != nil {
		req.AddCookie(c.CSRFCookie)
	}
	return req, nil
}

// OnRequestCompleted sets the DO API request completion callback
func (c *UniFiClient) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

func ConnectToSSHHost(user, host string) (*ssh.Client, *ssh.Session, error) {
	var pass string
	fmt.Print("Password: ")
	fmt.Scanf("%s\n", &pass)

	//var hostKey ssh.PublicKey
	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, nil, err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}
	defer session.Close()

	// Set IO
	session.Stdout = ansicolor.NewAnsiColorWriter(os.Stdout)
	session.Stderr = ansicolor.NewAnsiColorWriter(os.Stderr)
	in, _ := session.StdinPipe()

	// Set up terminal modes
	// https://net-ssh.github.io/net-ssh/classes/Net/SSH/Connection/Term.html
	// https://www.ietf.org/rfc/rfc4254.txt
	// https://godoc.org/golang.org/x/crypto/ssh
	// THIS IS THE TITLE
	// https://pythonhosted.org/ANSIColors-balises/ANSIColors.html
	modes := ssh.TerminalModes{
		ssh.ECHO:  0, // Disable echoing
		ssh.IGNCR: 1, // Ignore CR on input.
	}

	// Request pseudo terminal
	//if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
	//if err := session.RequestPty("xterm-256color", 80, 40, modes); err != nil {
	if err := session.RequestPty("vt100", 80, 40, modes); err != nil {
		//if err := session.RequestPty("vt220", 80, 40, modes); err != nil {
		log.Fatalf("request for pseudo terminal failed: %s", err)
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		log.Fatalf("failed to start shell: %s", err)
	}

	// Handle control + C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for {
			<-c
			fmt.Println("^C")
			fmt.Fprint(in, "\n")
			//fmt.Fprint(in, '\t')
		}
	}()

	//var b []byte = make([]byte, 1)

	// Accepting commands
	for {
		reader := bufio.NewReader(os.Stdin)
		str, _ := reader.ReadString('\n')
		fmt.Fprint(in, str)
	}

	return client, session, nil
}

// newResponse creates a new Response for the provided http.Response
func newResponse(r *http.Response) *Response {
	//dump, err := httputil.DumpResponse(r, true)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("%q", dump)

	response := Response{Response: r}

	//response.populateRate()

	return &response
}

func (r *Response) links() (map[string]headerLink.Link, error) {
	if linkText, ok := r.Response.Header["Link"]; ok {
		links, err := headerLink.Parse(linkText[0])

		if err != nil {
			return nil, err
		}

		linkMap := map[string]headerLink.Link{}
		for _, link := range links {
			linkMap[link.Rel] = link
		}

		return linkMap, nil
	}

	return map[string]headerLink.Link{}, nil
}

// populateRate parses the rate related headers and populates the response Rate.
func (r *Response) populateRate() {
	if limit := r.Header.Get(headerRateLimit); limit != "" {
		// r.Rate.Limit, _ = strconv.Atoi(limit)
	}
	if remaining := r.Header.Get(headerRateRemaining); remaining != "" {
		// r.Rate.Remaining, _ = strconv.Atoi(remaining)
	}
	if reset := r.Header.Get(headerRateReset); reset != "" {
		if v, _ := strconv.ParseInt(reset, 10, 64); v != 0 {
			// TODO: Fix this
			// r.Rate.Reset = Timestamp{time.Unix(v, 0)}
		}
	}
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *UniFiClient) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "unifises" {
			c.UnifiCookie = cookie
		}
		if cookie.Name == "csrf_token" {
			c.CSRFCookie = cookie
		}
	}

	response := newResponse(resp)
	//c.Rate = response.Rate

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, err
}

func (r *ErrorResponse) Error() string {
	if r.RequestID != "" {
		return fmt.Sprintf("%v %v: %d (request %q) %v",
			r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.RequestID, r.Message)
	}
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			return err
		}
	}

	return errorResponse
}

func (r Rate) String() string {
	return Stringify(r)
}

func (c *UniFiClient) Stop() {
	if c.Options.DbUsage.DbUsageEnabled {
		// Gracefully close database
		if err := c.Options.DbUsage.UnifiedDB.Close(); err != nil {
			panic(err)
		}
	}
}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}

// Int is a helper routine that allocates a new int32 value
// to store v and returns a pointer to it, but unlike Int32
// its argument value is an int.
func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

// StreamToString converts a reader to a string
func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(stream)
	return buf.String()
}
