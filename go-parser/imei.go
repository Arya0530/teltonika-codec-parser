package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func getPostgresURL() string {
	value := os.Getenv("POSTGRES_URL")
	if value != "" {
		return value
	}

	return "postgres://postgres:admin@" +
		"localhost:5434/iot_teltonika" +
		"?sslmode=disable"
}

func ValidateIMEI(imei string) (bool, error) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, getPostgresURL())
	if err != nil {
		return false, err
	}
	defer conn.Close(ctx)

	var exists bool

	err = conn.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM devices
			WHERE imei = $1
			  AND is_active = TRUE
		)
	`, imei).Scan(&exists)

	return exists, err
}

func CheckIMEIOrStop(imei string) bool {
	valid, err := ValidateIMEI(imei)
	if err != nil {
		fmt.Println("[IMEI] Database error:", err)
		return false
	}

	if !valid {
		fmt.Println("[IMEI] Device rejected:", imei)
		return false
	}

	fmt.Println("[IMEI] Device accepted:", imei)
	return true
}
