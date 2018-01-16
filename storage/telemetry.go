package storage

import (
	"time"
)

type Telemetry struct {
	Timestamp *time.Time `json:"time,omitempty"` // time=%TIME
	Provider  string     `json:"prov,omitempty"` // prov=%PROV

	Location Location `json:"location"`
	Device   struct {
		Battery   string `json:"battery,omitempty"`   // battery=%BATT
		Id        string `json:"deviceId,omitempty"`  //
		AndroidId string `json:"androidId,omitempty"` // androidId=%AID
		Serial    string `json:"serial,omitempty"`    // serial=%SER
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
