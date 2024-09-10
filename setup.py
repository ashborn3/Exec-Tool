import subprocess

# Define the commands
commands = [
    "del WindowsService.exe",
    "go build",
    "sc.exe stop mysvc",
    "sc.exe delete mysvc",
    "sc.exe create mysvc binPath=C:\\Users\\Nitin\\Desktop\\WindowsService\\WindowsService.exe",
    "sc.exe start mysvc"
]

# Execute each command
for command in commands:
    print(f"Executing: {command}")
    result = subprocess.run(command, shell=True, capture_output=True, text=True)
    
    if result.returncode != 0:
        print(f"Command failed: {command}")
        print(result.stderr)
        exit(1)
    else:
        print(result.stdout)

print("Service reloaded successfully with the latest Go build.")
