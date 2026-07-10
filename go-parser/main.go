package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

const maxAVLPacketSize = 1024 * 1024 // Maksimal 1 MB

func getTCPAddress() string {
	port := os.Getenv("TCP_PORT")
	if port == "" {
		port = "50020"
	}

	return ":" + port
}

func main() {
	address := getTCPAddress()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(fmt.Errorf("gagal menjalankan TCP server pada %s: %w", address, err))
	}
	defer listener.Close()

	fmt.Println("[TCP] Teltonika server listening on", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[TCP] Accept error:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("[TCP] Device connected:", conn.RemoteAddr())

	// Memberi waktu maksimal untuk proses handshake IMEI.
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Minute))

	imei, err := readIMEI(conn)
	if err != nil {
		fmt.Println("[IMEI] Read error:", err)
		return
	}

	fmt.Println("[IMEI] Received:", imei)

	if !CheckIMEIOrStop(imei) {
		if _, err := conn.Write([]byte{0x00}); err != nil {
			fmt.Println("[IMEI] Gagal mengirim ACK 00:", err)
		}

		fmt.Println("[IMEI] Rejected, ACK 00 sent")
		return
	}

	if _, err := conn.Write([]byte{0x01}); err != nil {
		fmt.Println("[IMEI] Gagal mengirim ACK 01:", err)
		return
	}

	fmt.Println("[IMEI] Accepted, ACK 01 sent")

	// Setelah handshake berhasil, tracker bisa mempertahankan koneksi.
	_ = conn.SetReadDeadline(time.Time{})

	for {
		hexPacket, err := readAVLPacket(conn)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				fmt.Println("[TCP] Device disconnected:", conn.RemoteAddr())
			} else {
				fmt.Println("[AVL] Read error:", err)
			}
			return
		}

		fmt.Println("[AVL] Packet HEX:", hexPacket)

		records, recordCount, err := DecodeAVLPacket(hexPacket, imei)
		if err != nil {
			fmt.Println("[AVL] Decode error:", err)
			return
		}

		// Setiap AVL record dipublish satu per satu ke RabbitMQ.
		for index, record := range records {
			jsonData, err := json.MarshalIndent(record, "", "  ")
			if err != nil {
				fmt.Println("[JSON] Marshal error:", err)
				return
			}

			fmt.Printf(
				"=== HASIL PARSER GOLANG RECORD %d/%d ===\n",
				index+1,
				recordCount,
			)
			fmt.Println(string(jsonData))

			if err := PublishRabbitMQ(record); err != nil {
				fmt.Println("[RabbitMQ] Error:", err)
				return
			}
		}

		// ACK 4 byte berisi jumlah record yang diterima.
		ack := make([]byte, 4)
		binary.BigEndian.PutUint32(ack, uint32(recordCount))

		if _, err := conn.Write(ack); err != nil {
			fmt.Println("[ACK] Write error:", err)
			return
		}

		fmt.Println("[ACK] Sent record count:", recordCount)
	}
}

func readIMEI(conn net.Conn) (string, error) {
	lengthBuf := make([]byte, 2)

	if _, err := io.ReadFull(conn, lengthBuf); err != nil {
		return "", err
	}

	imeiLength := binary.BigEndian.Uint16(lengthBuf)

	// IMEI Teltonika normalnya 15 karakter.
	if imeiLength == 0 || imeiLength > 32 {
		return "", fmt.Errorf("panjang IMEI tidak valid: %d", imeiLength)
	}

	imeiBuf := make([]byte, int(imeiLength))

	if _, err := io.ReadFull(conn, imeiBuf); err != nil {
		return "", err
	}

	return string(imeiBuf), nil
}

func readAVLPacket(conn net.Conn) (string, error) {
	// 4 byte preamble + 4 byte data length.
	header := make([]byte, 8)

	if _, err := io.ReadFull(conn, header); err != nil {
		return "", err
	}

	preamble := binary.BigEndian.Uint32(header[0:4])
	if preamble != 0 {
		return "", fmt.Errorf("invalid AVL preamble: %08X", preamble)
	}

	dataLength := binary.BigEndian.Uint32(header[4:8])

	if dataLength == 0 || dataLength > maxAVLPacketSize {
		return "", fmt.Errorf("AVL data length tidak valid: %d", dataLength)
	}

	// Data field + CRC 4 byte.
	body := make([]byte, int(dataLength)+4)

	if _, err := io.ReadFull(conn, body); err != nil {
		return "", err
	}

	fullPacket := append(header, body...)

	return hex.EncodeToString(fullPacket), nil
}
