package cloudcraft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httputil"
	"os"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

var apiKey = os.Getenv("CLOUDCRAFT_API_KEY")

func init() {
	log.SetLevel(log.DebugLevel)
	log.Debugf("CLOUDCRAFT_API_KEY: %s\n", apiKey)
}

const (
	baseURL   = "https://api.cloudcraft.co/"
	userAgent = "go-cloudcraft"
)

// A Client manages communication with the Cloudcraft API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for Cloudcraft requests.
	BaseURL *url.URL

	// User agent used when communicating with the Cloudcraft API.
	UserAgent string

	// Reuse a single struct instead of allocating one for each service on the heap.
	common service

	// Services used for talking to different parts of the Cloudcraft API.
	Blueprints *BlueprintsService
	// AWSAccounts *AWSAccountsService
	// Users       *UsersService
}

type service struct {
	client *Client
}

type Request http.Request

type Response http.Response

// An ErrorResponse reports one or more errors caused by an API request.
type ErrorResponse struct {
	// Response *http.Response
	Code    int    `json:"code"`  // status code
	Message string `json:"error"` // error message
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v", r.Code, r.Message)
}

// CONSIDER: http package doesn't recommend this approach:
//
// RoundTrip should not attempt to interpret the response. In
// particular, RoundTrip must return err == nil if it obtained
// a response, regardless of the response's HTTP status code.
// A non-nil err should be reserved for failure to obtain a
// response. Similarly, RoundTrip should not attempt to
// handle higher-level protocol details such as redirects,
// authentication, or cookies.
type transport struct {}

// RoundTrip wraps http.DefaultTransport.RoundTrip adding authentication header.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+apiKey)

	dump, _ := httputil.DumpRequestOut(req, true)
	log.Debugf("Request\n\n%s\n\n", dump)

	resp, err := http.DefaultTransport.RoundTrip(req)

	dump, _ = httputil.DumpResponse(resp, true)
	log.Debugf("Response\n\n%s\n\n", dump)

	return resp, err
}

// NewClient returns a new Cloudcraft API client. If a nil httpClient is
// provided, a new http.Client will be used. To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you.
func NewClient(httpClient *http.Client) *Client {
	t := &transport{}

	if httpClient == nil {
		httpClient = &http.Client{
			Transport: t,
		}
	}

	baseURL, _ := url.Parse(baseURL)

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}
	c.common.client = c
	c.Blueprints = (*BlueprintsService)(&c.common)

	return c
}

// NewRequest creates an API request. A relative URL can be provided in url,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, url string, body interface{}) (*Request, error) {
	u, err := c.BaseURL.Parse(url)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return (*Request)(req), nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v.
func (c *Client) Do(req *Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do((*http.Request)(req))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		errorResponse := &ErrorResponse{}
		data, err := ioutil.ReadAll(resp.Body)
		if err == nil && data != nil {
			json.Unmarshal(data, errorResponse)
		}
		return (*Response)(resp), errorResponse
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return nil, err
		}
	}

	return (*Response)(resp), err
}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }
