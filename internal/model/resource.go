package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Resource struct {
	Resource      JSONMap
	Cluster       string
	Timestamp     time.Time
	ResourceType  string
	ResourceName  string
	NamespaceName string
	AppName       string
	AppVersion    string
	GitHash       string
	State         string
}

type ResourceRequest struct {
	Application string    `json:"application"`
	Cluster     string    `json:"cluster"`
	GitHash     string    `json:"git_hash"`
	Time        time.Time `json:"time"`
	Type        string    `json:"type"`
	Version     string    `json:"version"`
	State       string    `json:"state"`
	Resources   []JSONMap `json:"resources"`
}

type JSONMap map[string]interface{}

func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	ba, err := m.MarshalJSON()
	return string(ba), err
}

func (m *JSONMap) Scan(val interface{}) error {
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return fmt.Errorf("Failed to unmarshal JSONB value:", val)
	}
	t := map[string]interface{}{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

func (m JSONMap) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	t := (map[string]interface{})(m)
	return json.Marshal(t)
}

func (m *JSONMap) UnmarshalJSON(b []byte) error {
	t := map[string]interface{}{}
	err := json.Unmarshal(b, &t)
	*m = t
	return err
}
