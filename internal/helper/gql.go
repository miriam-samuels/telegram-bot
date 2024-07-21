package helper

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// GraphQL query structure

type GraphQLRequest struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

// GraphQL response structure
type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []interface{}   `json:"errors"`
}

type Node struct {
	Nodes      []map[string]interface{} `json:"nodes"`
	TotalCount int                      `json:"totalCount"`
}

// GraphQL query structure
func FetchGraphQlData(reqBody *GraphQLRequest) (map[string]Node, error) {
	value, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Failed to marshal request body: %v", err)
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, os.Getenv("GRAPHQL_ENDPOINT"), bytes.NewBuffer(value))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 60}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("The HTTP response failed with error %s\n", err)
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Rrquest failed: %v", response)
		return nil, err
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, err
	}

	// Parse the response body
	var gqlResp GraphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		log.Printf("Failed to unmarshal response body: %v", err)
		return nil, err
	}

	// Check for errors in the GraphQL response
	if len(gqlResp.Errors) > 0 {
		log.Printf("GraphQL query errors: %v", gqlResp.Errors)
		return nil, err
	}

	var responseBody map[string]Node

	if err := json.Unmarshal(gqlResp.Data, &responseBody); err != nil {
		log.Printf("Failed to unmarshal user data: %v", err)
		return nil, err
	}

	return responseBody, nil
}
