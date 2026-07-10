def aggregate_data(data):
    speed = data.get("speed", 0)

    if speed == 0:
        status = "idle"
    elif speed > 60:
        status = "overspeed"
    else:
        status = "moving"

    data["status"] = status
    return data