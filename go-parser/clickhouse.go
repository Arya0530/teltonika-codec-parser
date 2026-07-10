package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func InsertClickHouse(data AVLData) error {
	raw, _ := json.Marshal(data)
	rawJSON := strings.ReplaceAll(string(raw), "'", "''")

	query := fmt.Sprintf(`
INSERT INTO raw_telemetry.avl_raw
(imei, timestamp, latitude, longitude, speed, raw_json, priority, altitude, angle, satellites)
VALUES ('%s', '%s', %f, %f, %d, '%s', %d, %d, %d, %d)
`,
		data.IMEI,
		data.Timestamp,
		data.Latitude,
		data.Longitude,
		data.Speed,
		rawJSON,
		data.Priority,
		data.Altitude,
		data.Angle,
		data.Satellites,
	)

	resp, err := http.Post(
		"http://arya:123456@localhost:8123/",
		"text/plain",
		bytes.NewBufferString(query),
	)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("clickhouse insert failed: %s", resp.Status)
	}

	fmt.Println("[ClickHouse] Raw data inserted")
	return nil
}
