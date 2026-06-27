import json

def parse_teltonika_coord(hex_str):
    unsigned_int = int(hex_str, 16)
    
    # Paksa jadi Signed 32-bit biar angka minus (Lintang Selatan) terbaca
    if unsigned_int & 0x80000000:
        signed_int = unsigned_int - 0x100000000
    else:
        signed_int = unsigned_int
        
    return signed_int / 10000000.0

# Contoh 4 Byte Hex
hex_longitude = "065B96E0" 
hex_latitude = "FDDC5730"

data = {
    "latitude": parse_teltonika_coord(hex_latitude),
    "longitude": parse_teltonika_coord(hex_longitude)
}

# Ubah jadi JSON
json_data = json.dumps(data, indent=2)

print("=== HASIL PARSER PYTHON ===")
print(json_data)