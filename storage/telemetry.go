package storage

import (
	"encoding/json"
	"strings"
	"time"
)

type Telemetry struct {
	Timestamp         *time.Time `json:"-"`              // time=%TIME
	OriginalTimestamp *time.Time `json:"time,omitempty"` // time=%TIME
	Provider          string     `json:"prov,omitempty"` // prov=%PROV

	Location Location `json:"location"`
	Device   struct {
		Id      string `json:"id,omitempty"`      //
		Battery string `json:"battery,omitempty"` // battery=%BATT
	} `json:"device"`
	Gps struct {
		Satellites string `json:"sat,omitempty"` // sat=%SAT
		Accuracy   string `json:"acc,omitempty"` // acc=%ACC
	} `json:"gps"`
	Vehicle struct {
		Speed     float64 `json:"spd,omitempty"` // spd=%SPD
		Direction float64 `json:"dir,omitempty"` // dir=%DIR
	} `json:"vehicle"`

	Annotation string `json:"desc,omitempty"` // desc=%DESC

	Activity string `json:"activity,omitempty"` // activity=%ACT
	// Epoch    string `json:"epoch,omitempty"`    // epoch=%TIMESTAMP"
}

type aliasTelemetry Telemetry

func (t Telemetry) Type() string {
	return "Telemetry"
}

func (t Telemetry) Pk() string {
	if t.Timestamp == nil {
		t.Timestamp = t.OriginalTimestamp
	}
	buff, _ := json.Marshal(t.Timestamp)
	//timestamp := time.Now().UTC().Format(time.RFC3339)
	timestamp := string(buff)
	return strings.Trim(timestamp, "\"")
}
