package main

import (
	"bytes"
	"embed"
	"encoding/json"
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

type SvcCommandRequest struct {
	Command string `json:"command" binding:"required"`
	Imp     string `json:"imp" binding:"required"`
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

		fmt.Println("Received command:", req.Command)

		cmd := exec.Command("cmd.exe", "/C", req.Command)
		output, err := cmd.Output()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"output": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"output": string(output)})
	})

	router.POST("/run-command-on-svc", func(c *gin.Context) {
		var req SvcCommandRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("Received command:", req.Command)

		url := "http://localhost:3232/execute"

		req = SvcCommandRequest{
			Command: req.Command,
			Imp:     req.Imp,
		}

		jsonData, err := json.Marshal(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling JSON"})
			return
		}

		httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
			return
		}
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending request"})
			return
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding response"})
			return
		}
		c.JSON(http.StatusOK, result)

	})

	router.Run(":31337")
}
