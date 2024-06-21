package api

import (
	"fmt"
	"time"

	"github.com/miriam-samuels/telegram-bot/internal/helper"
	"github.com/miriam-samuels/telegram-bot/internal/template"
)

func GetNftNews() (string, error) {
	reqData := helper.GraphQLRequest{
		Query: `
		query MyQuery($orderBy: String, $limit: Int) {
			nftNews(orderBy: $orderBy, limit: $limit) {
			  nodes {
				preview
				source
				title
				link
			  }
			}
		  }`,
		Variables: map[string]string{
			"limit":   "10",
			"orderBy": "publishDate",
		},
	}

	res, err := helper.FetchGraphQlData(&reqData)
	if err != nil {
		return "", fmt.Errorf("error pulling news data; %v", err)
	}

	message := helper.FormatHTMLMessage(res["nftNews"].Nodes, template.News)

	return message, nil
}

func GetSpaces() (string, error) {
	reqData := helper.GraphQLRequest{
		Query: `
		query MyQuery($orderBy: String, $scheduled: DateTime, $limit: Int) {
			twitterSpace(orderBy: $orderBy, scheduled: $scheduled, limit: $limit) {
			  nodes {
				title
				spaceUrl
				scheduled
				userhandle
			  }
			}
		  }`,
		Variables: map[string]string{
			"limit":   "10",
			"orderBy": "scheduled",
			"scheduled": time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		},
	}

	res, err := helper.FetchGraphQlData(&reqData)

	if err != nil {
		return "", fmt.Errorf("error pulling space data; %v", err)
	}

	message := helper.FormatHTMLMessage(res["twitterSpace"].Nodes, template.Spaces)

	return message, nil
}
