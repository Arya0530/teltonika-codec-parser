# Gunakan OS Linux versi ringan yang sudah ada Golang-nya
FROM golang:1.26-alpine

WORKDIR /app

COPY . .

# Build (compile) kode Go menjadi program aplikasi
RUN go build -o parser-app .

# Jalankan aplikasinya
CMD ["./parser-app"]