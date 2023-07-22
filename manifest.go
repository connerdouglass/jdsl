package jdsl

import (
	"encoding/json"
	"fmt"
	"os"
)

type Manifest struct {
	File      string   `json:"File"`
	Class     string   `json:"Class"`
	Author    string   `json:"Author"`
	Purpose   string   `json:"Purpose"`
	Functions []string `json:"Functions"`
}

func (m *Manifest) ReadFile(file string) error {
	// Read file contents
	contents, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("reading manifest file: %s", err)
	}

	// Unmarshal the JSON into the receiver
	if err := json.Unmarshal(contents, m); err != nil {
		return fmt.Errorf("unmarshaling manifest: %s", err)
	}
	return nil
}
