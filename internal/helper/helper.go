package helper

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func DivideFromSymbol(word string, symbol string) string {
	formattedWord, _, _ := strings.Cut(word, symbol)
	return formattedWord
}

func CheckWordSimilarity(word string, cmp string) bool {
	unspacedWord := strings.ToLower(strings.ReplaceAll(word, " ", ""))
	unspacedCmp := strings.ToLower(strings.ReplaceAll(cmp, " ", ""))

	// Return true if the words are already the same when unspaced or contains word
	if unspacedWord == unspacedCmp {
		return true
	}

	wordSlice := strings.Split(word, " ")
	cmpSlice := strings.Split(cmp, " ")

	//  only use this when the word sent by user is more than 1 eg: fox federation  wll give famous fox federation and transdimentional fox federation
	if strings.Contains(unspacedCmp, unspacedWord) && len(wordSlice) > 1 {
		return true
	}

	//  to avoid too mny option only use this when the word .... eg: gen2 will give smb gen2 ....
	if len(wordSlice) == 1 && len(cmpSlice) == 2 {
		for i := range wordSlice {
			for j := range cmpSlice {
				if strings.EqualFold(wordSlice[i], cmpSlice[j]) {
					return true
				}
			}
		}
	}

	// Create a copy of unspacedWord to modify during the loop
	numberMatched := characterMatch(unspacedWord, unspacedCmp)
	orderMatched := orderMatch(unspacedWord, unspacedCmp, 2)
	reverseOrderMatched := reverseOrderMatch(unspacedWord, unspacedCmp, 2)

	if reverseOrderMatched == len(unspacedWord) && len(wordSlice) > 1 {
		return true
	}

	if orderMatched == len(unspacedWord) && len(wordSlice) > 1 {
		return true
	}

	if len(unspacedWord) > 2 {
		return bool((numberMatched >= (len(unspacedWord)-2)) && (len(unspacedCmp) <= (len(unspacedWord)+1))) && orderMatched > 3
	}

	return false
}
func characterMatch(word string, cmp string) int {
	numberMatched := 0

	// Create a copy of unspacedWord to modify during the loop
	temp := word
	for _, ch := range cmp {
		if strings.ContainsRune(temp, ch) && len(temp) > 0 {
			numberMatched++
			temp = strings.Replace(temp, string(ch), "", 1)
		}
	}

	return numberMatched
}

func orderMatch(word string, cmp string, tolerance int) int {
	toleranceTaken := 0
	orderMatched := 0

	minimum := min(len(word), len(cmp))

	for i := 0; i < minimum; i++ {
		if word[i] == cmp[i] {
			orderMatched++
		} else {
			if toleranceTaken <= tolerance && i > 1 {
				orderMatched++
				toleranceTaken++
			} else {
				break
			}
		}
	}

	return orderMatched
}

func reverseOrderMatch(word string, cmp string, tolerance int) int {

	orderReverseMatched := 0

	toleranceTakenReverse := 0

	minimum := min(len(word), len(cmp))

	for i := 0; i < minimum; i++ {
		if word[len(word)-1-i] == cmp[len(cmp)-1-i] {
			orderReverseMatched++
		} else {
			if toleranceTakenReverse < tolerance && i < minimum-2 {
				orderReverseMatched++
				toleranceTakenReverse++
			} else {
				break
			}
		}
	}

	return orderReverseMatched
}

func TimeDiff(endDateStr string) int64 {
	// Parse the end date string into a time.Time object
	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		fmt.Println("Error parsing end date:", err)
		return 0
	}

	// Get the current date and time
	currentDate := time.Now().UTC()

	// Calculate the time difference
	timeDifference := endDate.Sub(currentDate).Seconds()

	// Check if the time difference is negative
	if timeDifference <= 0 {
		return 0
	}
	return int64(timeDifference)
}

func SendMessageToTelegram(chatID int, message, telegramToken string) error {
	// fmt.Println("sending message; ", message, chatID)
	var TelegramApiUrl = "https://api.telegram.org"
	url := fmt.Sprintf("%s/bot%s/sendMessage?chat_id=%d&text=%s", TelegramApiUrl, telegramToken, chatID, message)
	// fmt.Println("url; ", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed wrap request: %w", err)
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed send request: %w", err)
	}
	defer res.Body.Close()

	log.Printf("message sent successfully?\n%#v", res)

	return nil
}
