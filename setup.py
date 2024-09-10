import subprocess

# Define the commands
commands = [
    "del execService.exe",
    "go build src\\service.go",
    "move service.exe execService.exe",
    "sc.exe stop mysvc",
    "sc.exe delete mysvc",
    "sc.exe create mysvc binPath=C:/Users/Nitin/Desktop/Exec-Tool-main/execService.exe",
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
