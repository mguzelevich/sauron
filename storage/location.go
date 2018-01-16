package storage

import ()

type Location struct {
	Latitude  float64 `json:"lat"`           // lat=%LAT
	Longitude float64 `json:"lon"`           // lon=%LON
	Altitude  float64 `json:"alt,omitempty"` // alt=%ALT
}
