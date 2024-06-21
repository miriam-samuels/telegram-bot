package helper

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
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
	Nodes []DataItem `json:"nodes"`
}

type DataItem struct {
	Title      string `json:"title"`
	Preview    string `json:"preview"`
	Source     string `json:"source"`
	Link       string `json:"link"`
	Space      string `json:"spaceUrl"`
	UserHandle string `json:"userhandle"`
	Scheduled  string `json:"scheduled"`
}

// GraphQL query structure
func FetchGraphQlData(reqBody *GraphQLRequest) (map[string]Node, error) {
	value, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalf("Failed to marshal request body: %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, os.Getenv("GRAPHQL_ENDPOINT"), bytes.NewBuffer(value))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 60}

	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("The HTTP response failed with error %s\n", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Rrquest failed: %v", response)
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Parse the response body
	var gqlResp GraphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		log.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Check for errors in the GraphQL response
	if len(gqlResp.Errors) > 0 {
		log.Fatalf("GraphQL query errors: %v", gqlResp.Errors)
	}

	var responseBody map[string]Node

	if err := json.Unmarshal(gqlResp.Data, &responseBody); err != nil {
		log.Fatalf("Failed to unmarshal user data: %v", err)
	}

	return responseBody, nil
}

// Function to format data into HTML message
func FormatHTMLMessage(data []DataItem, tmpl string) string {
	funcMap := template.FuncMap{
		"add":        func(a, b int) int { return a + b },
		"capitalize": func(s string) string { return strings.ToUpper(string(s[0])) + s[1:] },
		"formatDate": func(t string) string {
			layout := "2006-01-02T15:04:05Z"
			d, err := time.Parse(layout, t)
			if err != nil {
				log.Printf("Error parsing date: %v", err)
				return t // return the original string if parsing fails
			}
			return d.Format("02 Jan - 15:04")
		},
		"cleanText": func(s string) string {
			words := strings.Fields(s)
			filteredWords := []string{}
			for _, word := range words {
				if !strings.HasPrefix(word, "#") {
					cleanedWord := strings.ReplaceAll(word, "@", "")
					filteredWords = append(filteredWords, cleanedWord)
				}
			}
			return strings.Join(filteredWords, " ")
		},
	}

	// Parse the template
	t, err := template.New("message").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}
	// cur retuned data
	var news []DataItem
	if len(data) < 10 {
		news = data
	} else {
		news = data[:10]
	}
	// Execute the template with the data
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, news); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	return tpl.String()
}
