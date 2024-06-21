import os

from dotenv import load_dotenv
import requests

load_dotenv()

enable_statistics = os.getenv("ENABLE_STATISTICS", "false").lower() == "true"
collector_url = os.getenv("STATISTICS_COLLECTOR_URL", "http://localhost:5000")


def load_worker_name():
    print("Loading worker name")
    if enable_statistics:
        response = requests.post(f"{collector_url}/worker-count")


def send_statistics(data):
    if enable_statistics:
        response = requests.post(f"{collector_url}/collect", json=data)
        print(response.json())
    else:
        print("Statistics disabled")
