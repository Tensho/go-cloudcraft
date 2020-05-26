package cloudcraft

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
)

// setup sets up a test HTTP server along with a cloudcraft.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)

	// client is the Cloudcraft client being tested and is configured to use test server.
	client = NewClient(nil)
	url, _ := url.Parse(server.URL + "/")
	log.Debugf("Server URL: %s", url)
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	type foo struct {
		bar string
	}

	inURL, outURL := "/blueprint", baseURL+"blueprint"
	inBody,  outBody := &Blueprint{Name: String("Web App")}, `{"name":"Web App"}`+"\n"
	req, _ := c.NewRequest("GET", inURL, inBody)

	// test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%v) Body is %v, want %v", inBody, got, want)
	}

	// test that default user-agent is attached to the request
	if got, want := req.Header.Get("User-Agent"), c.UserAgent; got != want {
		t.Errorf("NewRequest() User-Agent is %v, want %v", got, want)
	}
}

func TestDo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	type foo struct {
		Bar string `json:"bar"`
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"bar":"baz"}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	got := new(foo)
	client.Do(req, got)

	want := &foo{Bar: "baz"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response body = %v, want %v", got, want)
	}
}
