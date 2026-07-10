package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

type Codec8Reader struct {
	data []byte
	pos  int
}

func (r *Codec8Reader) read(n int) []byte {
	part := r.data[r.pos : r.pos+n]
	r.pos += n
	return part
}

func (r *Codec8Reader) readU8() uint8 {
	return r.read(1)[0]
}

func (r *Codec8Reader) readU16() uint16 {
	return binary.BigEndian.Uint16(r.read(2))
}

func (r *Codec8Reader) readU32() uint32 {
	return binary.BigEndian.Uint32(r.read(4))
}

func (r *Codec8Reader) readU64() uint64 {
	return binary.BigEndian.Uint64(r.read(8))
}

func coordToFloat(raw uint32) float64 {
	return float64(int32(raw)) / 10000000.0
}

func millisToTime(ms uint64) string {
	return time.UnixMilli(int64(ms)).Format("2006-01-02 15:04:05")
}

func skipCodec8IO(r *Codec8Reader) {
	r.readU8() // Event IO ID
	r.readU8() // Total IO

	n1 := r.readU8()
	for i := 0; i < int(n1); i++ {
		r.readU8()
		r.readU8()
	}

	n2 := r.readU8()
	for i := 0; i < int(n2); i++ {
		r.readU8()
		r.readU16()
	}

	n4 := r.readU8()
	for i := 0; i < int(n4); i++ {
		r.readU8()
		r.readU32()
	}

	n8 := r.readU8()
	for i := 0; i < int(n8); i++ {
		r.readU8()
		r.readU64()
	}
}

func skipCodec8EIO(r *Codec8Reader) {
	r.readU16() // Event IO ID
	r.readU16() // Total IO

	n1 := r.readU16()
	for i := 0; i < int(n1); i++ {
		r.readU16()
		r.readU8()
	}

	n2 := r.readU16()
	for i := 0; i < int(n2); i++ {
		r.readU16()
		r.readU16()
	}

	n4 := r.readU16()
	for i := 0; i < int(n4); i++ {
		r.readU16()
		r.readU32()
	}

	n8 := r.readU16()
	for i := 0; i < int(n8); i++ {
		r.readU16()
		r.readU64()
	}

	nx := r.readU16()
	for i := 0; i < int(nx); i++ {
		r.readU16()
		length := r.readU16()
		r.read(int(length))
	}
}

func DecodeAVLPacket(hexPacket string, imei string) ([]AVLData, int, error) {
	bytes, err := hex.DecodeString(hexPacket)
	if err != nil {
		return nil, 0, err
	}

	r := Codec8Reader{data: bytes}

	preamble := r.readU32()
	dataLength := r.readU32()
	codecID := r.readU8()
	recordCount := r.readU8()

	if preamble != 0 {
		return nil, 0, fmt.Errorf("invalid preamble")
	}

	if codecID != 0x08 && codecID != 0x8E {
		return nil, 0, fmt.Errorf("unsupported codec id: %X", codecID)
	}

	fmt.Println("Preamble:", preamble)
	fmt.Println("Data Length:", dataLength)
	fmt.Println("Codec ID:", fmt.Sprintf("%X", codecID))
	fmt.Println("Record Count:", recordCount)

	records := make([]AVLData, 0)

	for i := 0; i < int(recordCount); i++ {
		timestamp := r.readU64()
		priority := r.readU8()

		longitudeRaw := r.readU32()
		latitudeRaw := r.readU32()

		altitude := r.readU16()
		angle := r.readU16()
		satellites := r.readU8()
		speed := r.readU16()

		if codecID == 0x08 {
			skipCodec8IO(&r)
		} else if codecID == 0x8E {
			skipCodec8EIO(&r)
		}

		record := AVLData{
			IMEI:       imei,
			Timestamp:  millisToTime(timestamp),
			Latitude:   coordToFloat(latitudeRaw),
			Longitude:  coordToFloat(longitudeRaw),
			Speed:      speed,
			Priority:   priority,
			Altitude:   altitude,
			Angle:      angle,
			Satellites: satellites,
		}

		records = append(records, record)
	}

	// Number of Data 2
	totalRecords := r.readU8()
	fmt.Println("Total Records Footer:", totalRecords)

	// CRC 4 bytes
	r.readU32()

	return records, int(recordCount), nil
}
