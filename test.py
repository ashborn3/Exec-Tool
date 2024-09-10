import requests
import json

url = "http://localhost:3232/execute"

# Define the commands
commands = [
    "ipconfig",
    "ping google.com",
    "dir",
    "netstat",
    "tasklist",
    "systeminfo",
    "sfc /scannow",
    "chkdsk",
    "ipconfig /all",
]

# Set the headers
headers = {
    "Content-Type": "application/json"
}

# Iterate over the commands
for command in commands:
    # Define the payload
    payload = {
        "command": command,
        "imp": "system"
    }

    # Send the POST request with system integrity
    response = requests.post(url, headers=headers, data=json.dumps(payload))
    print(f"Command: {command}, Integrity: system")
    print("Status Code:", response.status_code)
    print("Response Body:", response.json())
    print()

    # Send the POST request with medium integrity
    payload["imp"] = "medium"
    response = requests.post(url, headers=headers, data=json.dumps(payload))
    print(f"Command: {command}, Integrity: medium")
    print("Status Code:", response.status_code)
    print("Response Body:", response.json())
    print()
    
    # Send the POST request with high integrity
    payload["imp"] = "high"
    response = requests.post(url, headers=headers, data=json.dumps(payload))
    print(f"Command: {command}, Integrity: high")
    print("Status Code:", response.status_code)
    print("Response Body:", response.json())
    print()

    # Send the POST request with low integrity (I know it doesn't exist. Just to show the error handling)
    payload["imp"] = "low"
    response = requests.post(url, headers=headers, data=json.dumps(payload))
    print(f"Command: {command}, Integrity: low")
    print("Status Code:", response.status_code)
    print("Response Body:", response.json())
    print()
