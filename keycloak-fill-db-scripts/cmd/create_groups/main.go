package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	baseUrl     = "http://localhost:8080"
	realm       = "test"
	accessToken = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJTd1ZEelA4VkhQdlBJbGRIQ3JjUWxOcDhmVlhpS1lyZDk3Z2IySGJQdXVVIn0.eyJleHAiOjE3MDEwNTM1NjcsImlhdCI6MTcwMTAxNzU2NywianRpIjoiMzAwYjdjYTktYzdmYy00NWNkLThhYWEtZjBjY2Q2ZjBkZDM3IiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL3JlYWxtcy9tYXN0ZXIiLCJzdWIiOiIzY2NjYTM2ZC05NjUxLTRmYjgtOGQzZC05YzdmZjZhZTQwYzciLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJhZG1pbi1jbGkiLCJzZXNzaW9uX3N0YXRlIjoiOWQ4MWQ2YmUtZGFjOS00ZTNhLTgxNTUtZGIxOWU3NzE3NzM4IiwiYWNyIjoiMSIsInNjb3BlIjoicHJvZmlsZSBlbWFpbCIsInNpZCI6IjlkODFkNmJlLWRhYzktNGUzYS04MTU1LWRiMTllNzcxNzczOCIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwicHJlZmVycmVkX3VzZXJuYW1lIjoidGVzdCJ9.Ih2DcyXzaYWw6jQweCmSV1jF5rYBugUWjBNWvQE9bCFLFO7J_wCtxHGg_LZGmHLQceI8WZVFR_SbyTH-cXrJPoCTRMrvHWGNFyYfjqsixKjXCjIkV8h7eZaxFrYPO035rTGb1j91rq41ooFYubRqviyyq5RwVy3dxlrjZlYgQnUf_gOuPgZdKWGQ7BY3FTdVv3iuxO1SWAxt1jFvGjcL4fJ_5opxjc0FxnDq6ksiYAmmXbc3lOvqTYB85DfOY35MXf5UvVCxguLIwmOpS4oLZY8zk9eMHZ6P-YcwHJs_kOnDPsv5aJaUXXg5UuE3Qk73xTC-Utme4WHf7oLkSkAocg"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run create_groups.go <numberOfGroups>")
		return
	}

	numberOfGroups, err := strconv.Atoi(os.Args[1])
	if err != nil || numberOfGroups <= 0 {
		fmt.Println("Please provide a valid positive integer for the number of groups.")
		return
	}

	client := &http.Client{}

	url := fmt.Sprintf("%s/admin/realms/%s/groups", baseUrl, realm)

	for i := 1; i <= numberOfGroups; i++ {
		groupName := fmt.Sprintf("group%d", i)
		payload := []byte(fmt.Sprintf(`{"name": "%s"}`, groupName))
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
		if err != nil {
			fmt.Printf("Error creating request for Group %d: %v\n", i, err)
			continue
		}

		req.Header.Set("Accept", "*/*")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request for Group %d: %v\n", i, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Error for Group %d. Status Code: %d, Response Body: %s\n", i, resp.StatusCode, body)
			return
		}

		fmt.Printf("Group %d created.\n", i)

		// Add a delay between requests to avoid overwhelming the server
		time.Sleep(5 * time.Millisecond)
	}
}
