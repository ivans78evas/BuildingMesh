package twenty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	apiURL string
	apiKey string
}

func NewClient(apiURL, apiKey string) *Client {
	return &Client{
		apiURL: apiURL,
		apiKey: apiKey,
	}
}

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

func (c *Client) Execute(query string, vars map[string]interface{}) (map[string]interface{}, error) {
	reqBody, _ := json.Marshal(GraphQLRequest{
		Query:     query,
		Variables: vars,
	})

	req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateProjectInTwenty creates a mirror of our project in Twenty for metadata management
func (c *Client) CreateProjectInTwenty(name, address string) (string, error) {
	query := `
		mutation CreateCompany($name: String!, $address: String) {
			createCompany(data: { name: $name, address: $address }) {
				id
			}
		}
	`
	vars := map[string]interface{}{
		"name":    name,
		"address": address,
	}

	res, err := c.Execute(query, vars)
	if err != nil {
		return "", err
	}

	// Parsing logic would go here in production
	fmt.Println("Twenty Response:", res)
	return "twenty-id-123", nil
}
