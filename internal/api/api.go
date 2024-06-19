package api

import (
	"fmt"
	"log"

	"github.com/miriam-samuels/telegram-bot/internal/helper"
)

func GetNftNews() string {
	reqData := helper.APIRequest{
		Method: "GET",
		Route:  "fetch-news",
	}

	res, err := helper.FetchData(&reqData)
	if err != nil {
		log.Fatalln("Error Occured")
	}

	// return fmt.Sprintf(
	// 	"Title: %s\nDescription: %s\nAttachment: %s\nLinks: %s\n",
	// 	user.Name,
	// 	user.Description,
	// 	user.Attachment,
	// 	user.Links,
	// )

	return res
}

func GetMintDrops() interface{} {
	queryString := `
	query ($limit: Int, $offset: Int) {
	mintCalendar(limit: $limit, offset: $offset) {
		nodes {
		  atlas3GiveawayUrl
		  blockchain
		  collectionId
		  createdAt
		  creatorEmail
		  creatorId
		  description
		  website
		  upvoteCount
		  updatedAt
		  twitterImageProfileUrl
		  twitterFollowers
		  twitter
		  thumbnailUrl
		  supply
		  subberPresaleUrl
		  subberGiveawayUrl
		  source
		  price
		  name
		  logoUrl
		  launchpadUrl
		  launchDate
		  isRejected
		  isOwn
		  isApproved
		  id
		  discord
		}
	  }
	  }`

	// Create a new GraphQL request
	reqBody := helper.GraphQLRequest{
		Query: queryString,
		Variables: map[string]interface{}{
			"limit":  "10",
			"offset": "10",
		},
	}

	res, err := helper.FetchGraphQlData(&reqBody)
	if err != nil {
		log.Fatalln("Error Occured")
	}

	fmt.Printf("ESP: %v", res)

	return res
}
