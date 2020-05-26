package cloudcraft

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBlueprintsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/blueprint", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"blueprints":[{"id":"6708d835-4c3a-40fe-a002-6103c8210510"}]}`)
	})

	got, _, err := client.Blueprints.List()
	if err != nil {
		t.Errorf("Blueprints.List() returned error: %v", err)
	}

	want := []*Blueprint{{ID: String("6708d835-4c3a-40fe-a002-6103c8210510")}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Blueprints.List() returned %+v, want %+v", got, want)
	}
}

func TestBlueprintsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/blueprint/6708d835-4c3a-40fe-a002-6103c8210510", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":"6708d835-4c3a-40fe-a002-6103c8210510"}`)
	})

	got, _, err := client.Blueprints.Get("6708d835-4c3a-40fe-a002-6103c8210510")
	if err != nil {
		t.Errorf("Blueprints.Get() returned error: %v", err)
	}

	want := &Blueprint{ID: String("6708d835-4c3a-40fe-a002-6103c8210510")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Blueprints.Get() returned %+v, want %+v", got, want)
	}
}
