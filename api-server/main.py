from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from postgres import get_tracking_data

app = FastAPI(title="Teltonika Tracking API")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
def home():
    return {
        "message": "Teltonika Tracking API is running"
    }

@app.get("/tracking")
def tracking():
    return get_tracking_data()