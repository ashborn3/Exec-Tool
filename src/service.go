package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

type postReqJson struct {
	Cmd string `json:"command"`
	Imp string `json:"imp"`
}

type Handler interface {
	Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (svcSpecificEC bool, exitCode uint32)
}

type myService struct{}

func (m *myService) Execute(
	args []string,
	r <-chan svc.ChangeRequest,
	status chan<- svc.Status) (svcSpecificEC bool,
	exitCode uint32,
) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	tick := time.Tick(30 * time.Second)
	status <- svc.Status{State: svc.StartPending}
	status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	go startServer()

loop:
	for {
		select {
		case <-tick:
			log.Print("Tick Handled...!", args)
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				status <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				log.Print("Shutting service...!")
				break loop
			case svc.Pause:
				status <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
				log.Printf("Unexpected service control request #%d", c)
			}
		}
	}

	status <- svc.Status{State: svc.StopPending}
	return false, 1
}

func runService(name string, isDebug bool) {
	if isDebug {
		err := debug.Run(name, &myService{})
		if err != nil {
			log.Fatalln("Error running service in debug mode.")
		}
	} else {
		err := svc.Run(name, &myService{})
		if err != nil {
			log.Fatalln("Error running service in Service Control mode.")
		}
	}
}

func startServer() {
	router := gin.Default()

	router.POST("/execute", func(ctx *gin.Context) {
		var postJson postReqJson
		err := ctx.BindJSON(&postJson)
		if err != nil {
			log.Println("Json Binding Failed!", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		log.Println("Command:", postJson.Cmd)
		log.Println("Impersonation:", postJson.Imp)

		// Medium Integrity: CreateProcessAsUser() after obtaining a medium integrity access token
		// High Integrity: CreateProcessAsUser() after obtaining a high-integrity access token
		// System Integrity: Simply use CreateProcess()
		var response string
		switch postJson.Imp {
		case "medium":
			response = executeWithMedInt(postJson.Cmd)
		case "high":
			response = executeWithHighInt(postJson.Cmd)
		case "system":
			response = runWithSystemInt(postJson.Cmd)
		default:
			log.Println("Invalid Impersonation Level!")
			response = "Invalid Impersonation Level!"
		}

		ctx.JSON(http.StatusOK, gin.H{"output": response})
	})

	router.Run(":3232")
}

func executeWithMedInt(cmd string) string {
	process := "cmd.exe"
	targetCmdLine := windows.StringToUTF16Ptr(process)

	var startupInfo windows.StartupInfo
	var processInfo windows.ProcessInformation

	err := windows.CreateProcess(nil, targetCmdLine, nil, nil, false, windows.CREATE_NEW_CONSOLE, nil, nil, &startupInfo, &processInfo)
	if err != nil {
		log.Fatalf("Failed to start process %s: %v", process, err)
	}
	defer windows.CloseHandle(processInfo.Process)
	defer windows.CloseHandle(processInfo.Thread)

	var token windows.Token
	err = windows.OpenProcessToken(processInfo.Process, windows.TOKEN_QUERY|windows.TOKEN_DUPLICATE, &token)
	if err != nil {
		log.Fatalf("Failed to get token of process %s: %v", process, err)
		return err.Error()
	}

	defer token.Close()
	var duplicatedToken windows.Token
	err = windows.DuplicateTokenEx(token, windows.TOKEN_ALL_ACCESS, nil, windows.SecurityImpersonation, windows.TokenPrimary, &duplicatedToken)
	if err != nil {
		log.Fatalf("Failed to duplicate token: %v", err)
	}
	defer duplicatedToken.Close()

	command := exec.Command("cmd.exe", "/C", cmd)

	command.SysProcAttr = &syscall.SysProcAttr{
		Token: syscall.Token(duplicatedToken),
	}

	output, err := command.Output()
	if err != nil {
		log.Fatalf("Failed to start process %s: %v", cmd, err)
	}

	log.Println(string(output))

	return string(output)
}

func executeWithHighInt(cmd string) string {
	process := "svchost.exe"
	targetCmdLine := windows.StringToUTF16Ptr(process)

	var startupInfo windows.StartupInfo
	var processInfo windows.ProcessInformation

	err := windows.CreateProcess(nil, targetCmdLine, nil, nil, false, windows.CREATE_NEW_CONSOLE, nil, nil, &startupInfo, &processInfo)
	if err != nil {
		log.Fatalf("Failed to start process %s: %v", process, err)
	}
	defer windows.CloseHandle(processInfo.Process)
	defer windows.CloseHandle(processInfo.Thread)

	var token windows.Token
	err = windows.OpenProcessToken(processInfo.Process, windows.TOKEN_QUERY|windows.TOKEN_DUPLICATE, &token)
	if err != nil {
		log.Fatalf("Failed to get token of process %s: %v", process, err)
		return err.Error()
	}

	defer token.Close()
	var duplicatedToken windows.Token
	err = windows.DuplicateTokenEx(token, windows.TOKEN_ALL_ACCESS, nil, windows.SecurityImpersonation, windows.TokenPrimary, &duplicatedToken)
	if err != nil {
		log.Fatalf("Failed to duplicate token: %v", err)
	}
	defer duplicatedToken.Close()

	command := exec.Command("cmd.exe", "/C", cmd)

	command.SysProcAttr = &syscall.SysProcAttr{
		Token: syscall.Token(duplicatedToken),
	}

	output, err := command.Output()
	if err != nil {
		log.Fatalf("Failed to start process %s: %v", cmd, err)
	}

	log.Println(string(output))
	return string(output)
}

func runWithSystemInt(cmd string) string {
	command := exec.Command("cmd.exe", "/C", cmd)
	output, err := command.Output()
	if err != nil {
		log.Fatalf("Failed to start process %s: %v", cmd, err)
	}

	log.Println(string(output))
	return string(output)
}

func main() {
	f, err := os.OpenFile("C:/Debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(fmt.Errorf("error opening file: %v", err))
	}
	defer f.Close()

	log.SetOutput(f)
	runService("myservice", false)
}
