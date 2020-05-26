package cloudcraft

import (
	"fmt"

	"github.com/hackebrot/go-repr/repr"
)

// BlueprintsService handles communication with the blueprint related
// methods of the Cloudcraft API.
//
// Cloudcraft API docs: https://developers.cloudcraft.co/
type BlueprintsService service

type Root struct {
	Blueprints []*Blueprint `json:"blueprints"`
}

// Blueprint represents a Cloudcraft blueprint.
type Blueprint struct {
	ID          *string   `json:"id,omitempty"`
	Name        *string   `json:"name,omitempty"`
	CreatedAt   *string   `json:"createdAt,omitempty"`
	UpdatedAt   *string   `json:"updatedAt,omitempty"`
	ReadAccess  []*string `json:"readAccess,omitempty"`
	WriteAccess []*string `json:"writeAccess,omitempty"`
	CreatorID   *string   `json:"CreatorId,omitempty"`
	LastUserID  *string   `json:"LastUserId,omitempty"`
}

func (b Blueprint) String() string {
	return repr.Repr(&b)
}

// List blueprints for a user.
//
// Cloudcraft API docs: https://developers.cloudcraft.co/
func (s *BlueprintsService) List() ([]*Blueprint, *Response, error) {
	req, err := s.client.NewRequest("GET", "blueprint", nil)
	if err != nil {
		return nil, nil, err
	}

	root := Root{}
	resp, err := s.client.Do(req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root.Blueprints, resp, nil
}

// Get retrieves blueprint by ID.
//
// Cloudcraft API docs: https://developers.cloudcraft.co/
func (s *BlueprintsService) Get(id string) (*Blueprint, *Response, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("blueprint/%s", id), nil)
	if err != nil {
		return nil, nil, err
	}

	blueprint := new(Blueprint)
	resp, err := s.client.Do(req, blueprint)
	if err != nil {
		return nil, resp, err
	}

	return blueprint, resp, nil
}
