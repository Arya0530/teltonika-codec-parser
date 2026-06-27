# Gunakan OS Linux versi ringan yang sudah ada Python-nya
FROM python:3.9-alpine

# Bikin folder kerja di dalam container
WORKDIR /app

# Copy semua file dari folder laptop lo ke dalam container
COPY . .

# Perintah otomatis saat container jalan
CMD ["python", "main.py"]