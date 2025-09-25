package main

import (
	"archive/zip"
	"encoding/json"
)

type Journal struct {
	Metadata Metadata `json:"metadata"`
	Entries  []*Entry `json:"entries"`
}

type Entry struct {
	RichText            string   `json:"richText"`
	Location            struct{} `json:"location"`
	TimeZone            string   `json:"timeZone"`
	IsPinned            bool     `json:"isPinned"`
	Text                string   `json:"text"`
	IsAllDay            bool     `json:"isAllDay"`
	CreationOSVersion   string   `json:"creationOSVersion"`
	ModifiedDate        string   `json:"modifiedDate"`
	Starred             bool     `json:"starred"`
	CreationDate        string   `json:"creationDate"`
	Tags                []string `json:"tags"`
	Weather             struct{} `json:"weather"`
	EditingTime         float64  `json:"editingTime"`
	UUID                string   `json:"uuid"`
	Duration            int      `json:"duration"`
	CreationDevice      string   `json:"creationDevice"`
	Photos              []*Photo `json:"photos"`
	CreationDeviceModel string   `json:"creationDeviceModel"`
	CreationDeviceType  string   `json:"creationDeviceType"`
	CreationOSName      string   `json:"creationOSName"`
}

type Metadata struct {
	Version string `json:"version"`
}

func parseDOJson(f *zip.File) (*Journal, error) {
	zfile, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer zfile.Close()

	dec := json.NewDecoder(zfile)
	j := &Journal{}
	err = dec.Decode(j)
	if err != nil {
		return nil, err
	}

	return j, nil
}
