package main

import (
	"embed"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

//go:embed frontend
var content embed.FS

type CommandRequest struct {
	Command string `json:"command" binding:"required"`
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		data, _ := content.ReadFile("frontend/index.html")
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	router.POST("/run-command", func(c *gin.Context) {
		var req CommandRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Process the command (for now, just echo it back and print it to console)
		fmt.Println("Received command:", req.Command)

		cmd := exec.Command(req.Command)
		output, err := cmd.Output()
		fmt.Println("Command output:", output)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"output": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"output": string(output)})
	})

	router.Run(":31337")
}
