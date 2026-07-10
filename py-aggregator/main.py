import json
import consumer

def parse_teltonika_coord(hex_str):
    unsigned_int = int(hex_str, 16)
    
    if unsigned_int & 0x80000000:
        signed_int = unsigned_int - 0x100000000
    else:
        signed_int = unsigned_int
        
    return signed_int / 10000000.0

# 4 Byte Hex
hex_longitude = "065B96E0" 
hex_latitude = "FDDC5730"

data = {
    "latitude": parse_teltonika_coord(hex_latitude),
    "longitude": parse_teltonika_coord(hex_longitude)
}

# JSON
json_data = json.dumps(data, indent=2)

print("=== HASIL PARSER PYTHON ===")
print(json_data)