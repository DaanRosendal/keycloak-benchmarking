package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseUrl     = "http://localhost:8080"
	realm       = "benchmarking"
	accessToken = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJTd1ZEelA4VkhQdlBJbGRIQ3JjUWxOcDhmVlhpS1lyZDk3Z2IySGJQdXVVIn0.eyJleHAiOjE3MDEwNTM1NjcsImlhdCI6MTcwMTAxNzU2NywianRpIjoiMzAwYjdjYTktYzdmYy00NWNkLThhYWEtZjBjY2Q2ZjBkZDM3IiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL3JlYWxtcy9tYXN0ZXIiLCJzdWIiOiIzY2NjYTM2ZC05NjUxLTRmYjgtOGQzZC05YzdmZjZhZTQwYzciLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJhZG1pbi1jbGkiLCJzZXNzaW9uX3N0YXRlIjoiOWQ4MWQ2YmUtZGFjOS00ZTNhLTgxNTUtZGIxOWU3NzE3NzM4IiwiYWNyIjoiMSIsInNjb3BlIjoicHJvZmlsZSBlbWFpbCIsInNpZCI6IjlkODFkNmJlLWRhYzktNGUzYS04MTU1LWRiMTllNzcxNzczOCIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwicHJlZmVycmVkX3VzZXJuYW1lIjoidGVzdCJ9.Ih2DcyXzaYWw6jQweCmSV1jF5rYBugUWjBNWvQE9bCFLFO7J_wCtxHGg_LZGmHLQceI8WZVFR_SbyTH-cXrJPoCTRMrvHWGNFyYfjqsixKjXCjIkV8h7eZaxFrYPO035rTGb1j91rq41ooFYubRqviyyq5RwVy3dxlrjZlYgQnUf_gOuPgZdKWGQ7BY3FTdVv3iuxO1SWAxt1jFvGjcL4fJ_5opxjc0FxnDq6ksiYAmmXbc3lOvqTYB85DfOY35MXf5UvVCxguLIwmOpS4oLZY8zk9eMHZ6P-YcwHJs_kOnDPsv5aJaUXXg5UuE3Qk73xTC-Utme4WHf7oLkSkAocg"
	userID      = "3a0dcdbf-e19c-424c-b57c-efbc3757a917"
)

type Group struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Path      string   `json:"path"`
	SubGroups []string `json:"subGroups"`
}

func main() {
	// Step 1: Fetch groups
	groups, err := fetchGroups()
	if err != nil {
		fmt.Println("Error fetching groups:", err)
		return
	}

	fmt.Printf("Fetched %d groups\n", len(groups))

	// Step 2: Add user to each group
	for _, group := range groups {
		err := addUserToGroup(group.ID)
		if err != nil {
			fmt.Printf("Error adding user to group %s: %v\n", group.Name, err)
		} else {
			fmt.Printf("User added to group: %s\n", group.Name)
		}

		// Add a delay between requests to avoid overwhelming the server
		time.Sleep(25 * time.Millisecond)
	}

	fmt.Println("User added to all groups")
}

func fetchGroups() ([]Group, error) {
	url := fmt.Sprintf("%s/admin/realms/%s/groups", baseUrl, realm)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	fmt.Println("Fetching groups...")

	res, err := http.DefaultClient.Do(req)
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

func addUserToGroup(groupID string) error {
	url := fmt.Sprintf("%s/admin/realms/%s/users/%s/groups/%s", baseUrl, realm, userID, groupID)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check the response status code
	if res.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to add user to group. Status: %d, Response: %s", res.StatusCode, body)
	}

	return nil
}
