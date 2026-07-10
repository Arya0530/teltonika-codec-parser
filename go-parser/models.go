package main

type AVLData struct {
	IMEI       string  `json:"imei"`
	Timestamp  string  `json:"timestamp"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Speed      uint16  `json:"speed"`
	Priority   uint8   `json:"priority"`
	Altitude   uint16  `json:"altitude"`
	Angle      uint16  `json:"angle"`
	Satellites uint8   `json:"satellites"`
}
