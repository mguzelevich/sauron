package sauron

import (
//
)

type Location struct {
	Lat       string `json:"lat"`                 // lat=%LAT
	Lon       string `json:"lon"`                 // lon=%LON
	Sat       string `json:"sat,omitempty"`       // sat=%SAT
	Desc      string `json:"desc,omitempty"`      // desc=%DESC
	Alt       string `json:"alt,omitempty"`       // alt=%ALT
	Acc       string `json:"acc,omitempty"`       // acc=%ACC
	Dir       string `json:"dir,omitempty"`       // dir=%DIR
	Prov      string `json:"prov,omitempty"`      // prov=%PROV
	Spd       string `json:"spd,omitempty"`       // spd=%SPD
	Time      string `json:"time,omitempty"`      // time=%TIME
	Battery   string `json:"battery,omitempty"`   // battery=%BATT
	AndroidId string `json:"androidId,omitempty"` // androidId=%AID
	Serial    string `json:"serial,omitempty"`    // serial=%SER
	Activity  string `json:"activity,omitempty"`  // activity=%ACT
	Epoch     string `json:"epoch,omitempty"`     // epoch=%TIMESTAMP"
}
