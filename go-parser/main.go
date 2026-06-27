package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Struct kerangka data JSON
type AVLData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func main() {
	// Contoh 4 Byte Hex dari GPS
	hexLongitude := "065B96E0" 
	hexLatitude := "FDDC5730"  

	// Panggil fungsi pemotong hex
	lon := parseTeltonikaCoord(hexLongitude)
	lat := parseTeltonikaCoord(hexLatitude)

	data := AVLData{
		Latitude:  lat,
		Longitude: lon,
	}

	// Ubah jadi JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("=== HASIL PARSER GOLANG ===")
	fmt.Println(string(jsonData))
}

// Fungsi inti ubah Hex ke Desimal
func parseTeltonikaCoord(hexStr string) float64 {
	parsedUint, _ := strconv.ParseUint(hexStr, 16, 32)
	signedInt := int32(parsedUint) // Paksa jadi minus kalau Lintang Selatan
	return float64(signedInt) / 10000000.0
}