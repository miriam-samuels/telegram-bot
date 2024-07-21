package helper

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/dustin/go-humanize"
)

var funcMap = template.FuncMap{
	"add":          add,
	"minus":        minus,
	"divide":       divide,
	"multiply":     multiply,
	"capitalize":   capitalize,
	"formatDate":   formatDate,
	"cleanText":    formatText,
	"isLive":       isTimePassed,
	"greater":      greater,
	"formatLink":   FormatLink,
	"addComma":     addComma,
	"poolDuration": poolDuration,
	"convertTime":  convertTime,
	"timeTo":       timeTo,
}

func add(a, b int) int           { return a + b }
func minus(a, b float64) float64 { return a - b }

func divide(a, b float64) string {
	d := a / b
	formatted := fmt.Sprintf("%.2f", d)

	parts := strings.Split(formatted, ".")

	val, _ := strconv.ParseInt(parts[0], 10, 64)
	return humanize.Comma(val) + "." + parts[1]
}

func multiply(a, b float64) string {
	d := a * b
	formatted := fmt.Sprintf("%.2f", d)
	parts := strings.Split(formatted, ".")

	val, _ := strconv.ParseInt(parts[0], 10, 64)
	return humanize.Comma(val) + "." + parts[1]
}
func addComma(val interface{}) string {
	switch val.(type) {
	case int:
		return humanize.Comma(int64(val.(int)))
	case float64:
		return humanize.Commaf(val.(float64))
	default:
		return humanize.Comma(val.(int64))

	}
}

func greater(a, b float64) bool {
	if a > b {
		return true
	} else {
		return false
	}
}

func capitalize(s string) string { return strings.ToUpper(string(s[0])) + s[1:] }

func formatDate(t string) string {
	layout := "2006-01-02T15:04:05Z"
	d, err := time.Parse(layout, t)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		return t // return the original string if parsing fails
	}
	return d.Format("02 Jan - 15:04")
}

func formatText(s string) string {
	words := strings.Fields(s)
	filteredWords := []string{}
	for _, word := range words {
		if !strings.HasPrefix(word, "#") {
			cleanedWord := strings.ReplaceAll(word, "@", "")
			filteredWords = append(filteredWords, cleanedWord)
		}
	}
	return strings.Join(filteredWords, " ")
}

func isTimePassed(timeStr string) (bool, error) {
	// Define the layout for the time string
	layout := "2006-01-02T15:04:05Z"

	// Parse the given time string to time.Time
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return false, fmt.Errorf("error parsing time: %v", err)
	}

	// Get the current time in UTC
	currentTime := time.Now().UTC()

	// Compare the parsed time with the current UTC time
	return parsedTime.Before(currentTime), nil
}

// ConvertTime converts a duration given in seconds to days
// func convertTime(seconds int64) string {
// 	days := seconds / (24 * 60 * 60)
// 	return fmt.Sprintf("%d days", days)
// }

func convertTime(seconds int64) string {
	var timeRemaining string
	days := seconds / (24 * 60 * 60)
	if days > 1 {
		timeRemaining = timeRemaining + fmt.Sprintf("%d", days) + "days "
	}

	hours := (seconds % (24 * 60 * 60)) / (60 * 60)
	if hours > 1 {
		timeRemaining = timeRemaining + fmt.Sprintf("%d", hours) + "hrs "
	}

	minutes := (seconds % (60 * 60)) / 60
	if minutes > 1 {
		timeRemaining = timeRemaining + fmt.Sprintf("%d", minutes) + "min "
	}

	sec := (seconds % (60)) / 60
	if sec > 60 {
		timeRemaining = timeRemaining + fmt.Sprintf("%d", seconds) + "sec "
	}
	return timeRemaining
}

func poolDuration(pool map[string]interface{}) string {
	marketplace := strings.ToLower(pool["marketplace"].(string))
	if marketplace == "citrus" {
		return "7 - 14 days"
	} else if marketplace == "banx" && pool["collectionName"].(string) == "Flip loans (Pool)" {
		return "7/14 days"
	} else if marketplace == "banx" {
		return "Perpetual"
	} else {
		return convertTime(int64(pool["duration"].(float64)))
	}
}

func timeTo(endDateStr string) string {
	seconds := TimeDiff(endDateStr)

	return convertTime(seconds)
}

func FormatLink(s string) string { return strings.ReplaceAll(s, " ", "-") }

// Function to format data into HTML message
func FormatHTMLMessage(data []map[string]interface{}, tmpl string) (string, error) {
	// Parse the template
	t, err := template.New("message").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		return "", err
	}

	// Execute the template with the data
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		log.Printf("Error executing template: %v", err)
		return "", err
	}

	return tpl.String(), nil
}
