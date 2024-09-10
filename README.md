# Exec-Tool

Exec-Tool is a tool to easily execute commands, run scripts, and automate repetitive tasks, all from the comfort of your browser.

## Build Instructions

To build Exec-Tool, follow these steps:

1. Clone the repository: `git clone https://github.com/your-username/Exec-Tool.git`
2. Navigate to the project directory: `cd Exec-Tool`
3. Build the tool: `go build src/main.go`
4. Start using Exec-Tool: `./main`

## Service Installation

To install the service, follow these steps:

1. After Cloning the repo and building the tool as sin the shown in the Build Instructions Section 
2. Build the service using `go build src/service.go`
3. In the powershell, run `sc.exe create <service name> binPath=<absolute path to the service executable file>`
4. run `sc.exe start <service name>`
5. To delete the service, run `sc.exe stop <service name>` followed by `sc.exe delete <service name>` 

## Usage

- Run the tool: `./main`
- Execute a command: Enter the Command in the first box and press enter or click on Run Command
- Output will be displayed in the box below the button

