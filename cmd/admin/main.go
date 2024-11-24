package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
)

func init() {
	initCmd()
}

type App struct {
	Host string
}

var app *App

var cmdRoot = &cobra.Command{
	Use: "rock enterprise backup admin utilities",
}

var jobAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new job",
	Run: func(cmd *cobra.Command, args []string) {
		jobType, _ := cmd.Flags().GetString("type")
		backupType, _ := cmd.Flags().GetString("backup_type")
		policy_id, _ := cmd.Flags().GetString("policy_id")

		if jobType == "backup" {
			app.StartBackup(policy_id, backupType)
		}
	},
}

func initCmd() {
	jobCmd.AddCommand(jobAddCmd)
	cmdRoot.AddCommand(jobCmd)
}

// main admin job start backup
var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "job utilities",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func main() {

	app = &App{"http://localhost:8000"}

	jobAddCmd.Flags().StringP("type", "t", "", "Type of the job (e.g., backup)")
	jobAddCmd.Flags().StringP("backup_type", "b", "full", "Data type of the job (e.g., full, incremental)")
	jobAddCmd.Flags().StringP("policy_id", "p", "0", "id of policy")

	if err := cmdRoot.Execute(); err != nil {
		panic(err)
	}
}

func (a *App) StartBackup(policyID, backupType string) {
	// url := "backup/job"

	baseURL, err := url.Parse(a.Host)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	baseURL.Path += "/backup/job"

	httpClient := &http.Client{}

	data := map[string]string{"policy_id": policyID, "backup_type": backupType}
	fmt.Printf("%v\n\n", data)
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", baseURL.String(), bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	resp, err := httpClient.Do(req)

	if err != nil {
		fmt.Println("http client error:", err)
		return
	}

	fmt.Printf("%v\n", resp)
}
