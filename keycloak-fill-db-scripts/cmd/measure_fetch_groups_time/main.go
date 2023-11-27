package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

const (
	baseUrl     = "http://localhost:8080"
	realm       = "benchmarking"
	accessToken = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJtbllHSlNNVGE4bXh0dzhtUElIeWgxVWViVDdWYTQycFBwWFlwWGpYS2xvIn0.eyJleHAiOjE3MDExMzQ2NzksImlhdCI6MTcwMTA5ODY3OSwianRpIjoiNTQxYmM4ZTUtMTEzMi00ZTQxLWI3ODUtZjAwNDRmZWE2NDFhIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL3JlYWxtcy9tYXN0ZXIiLCJzdWIiOiJhZTBmZTE5ZS1mNjFiLTQyNWItODUyOC1mNzU3MGYzNzMzZGQiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJhZG1pbi1jbGkiLCJzZXNzaW9uX3N0YXRlIjoiMDRlYmQ4MGUtYWQxNy00ZTc5LTg3YjUtM2IyZDA1Y2M0YzE5IiwiYWNyIjoiMSIsInNjb3BlIjoicHJvZmlsZSBlbWFpbCIsInNpZCI6IjA0ZWJkODBlLWFkMTctNGU3OS04N2I1LTNiMmQwNWNjNGMxOSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwicHJlZmVycmVkX3VzZXJuYW1lIjoiYWRtaW4ifQ.ke9QlkeS3zl3zVFuJrVHmGVgkDciB5GgCol3I_M-NDeQpuPuAZ41vN7M_ZwkENGQ86JgRuExC3PVmLcxstkLG-ebz8CYlw28o8ClzEmG-d2fVn_ZXqhuLw6VFqWIMLBkFn697tlFq6UGxZIUiJtNbFKa02kHhn96SMhK10rqVJTq3RKyJE79_7-u42HwQB06mEEf9TCvHe1zSJZSIP0H9VYoTxeqtSjU0PpwA-eui096qxWbqEVZNN2dTpMVFDoBLZCxXvAaHCkoecBkM_Upmleu12RPRZthtNEtmVJqL_USZ9Kmrmz1zE16q2BBl3E0Dzi9snqqLO_mLU61U1PZmA"
)

type Group struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Path      string   `json:"path"`
	SubGroups []string `json:"subGroups"`
}

var createdGroupsCounter int

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run measure_fetch_groups_time.go <maxExponent>")
		return
	}

	maxExponent, err := strconv.Atoi(os.Args[1])
	if err != nil || maxExponent <= 0 {
		fmt.Println("Please provide a valid positive integer for the max exponent.")
		return
	}

	client := &http.Client{}
	url := fmt.Sprintf("%s/admin/realms/%s/groups", baseUrl, realm)

	// Create a new table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Iteration", "Number of Groups", "Fetch Time"})

	createdGroupsCounter = 0

	for i := 1; i <= maxExponent; i++ {
		numGroups := int(math.Pow(2, float64(i)))

		err := createGroups(client, url, createdGroupsCounter, numGroups)
		if err != nil {
			fmt.Println("Error creating groups:", err)
			return
		}

		fetchTime := measureFetchTime(client, url)

		// Add data to the table
		table.Append([]string{strconv.Itoa(i), strconv.Itoa(createdGroupsCounter), fetchTime.String()})

		fmt.Printf("Iteration %d - Number of Groups: %d, Fetch Time: %v\n\n", i, numGroups, fetchTime)
	}

	// Render the table
	table.Render()
}

func createGroups(client *http.Client, url string, startIndex, numGroups int) error {

	fmt.Printf("Creating %d groups...\n", numGroups)

	for i := startIndex; i < numGroups; i++ {
		createdGroupsCounter++
		groupName := fmt.Sprintf("group%d", startIndex+i)
		payload := []byte(fmt.Sprintf(`{"name": "%s"}`, groupName))
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
		if err != nil {
			fmt.Printf("Error creating request for Group %d: %v\n", i, err)
			return err
		}

		req.Header.Set("Accept", "*/*")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request for Group %d: %v\n", i, err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			fmt.Printf("Error for Group %d. Status Code: %d\n", i, resp.StatusCode)
			return errors.New("failed to create group")
		}

		// Add a delay between requests to avoid overwhelming the server
		time.Sleep(5 * time.Millisecond)
	}

	fmt.Printf("Created %d groups.\n", numGroups)

	return nil
}

func measureFetchTime(client *http.Client, url string) time.Duration {
	startTime := time.Now()

	_, err := fetchGroups(client, url)
	if err != nil {
		fmt.Println("Error fetching groups:", err)
	}

	endTime := time.Now()
	return endTime.Sub(startTime)
}

func fetchGroups(client *http.Client, url string) ([]Group, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	fmt.Printf("Fetching %d groups...\n", createdGroupsCounter)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch groups. Status: %d, Response: %s", res.StatusCode, body)
	}

	var groups []Group
	err = json.Unmarshal(body, &groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}
