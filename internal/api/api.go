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
	// Define the HTML template
	const tmpl = `
<b>Latest NFT News (from <a href="https://pr-1540.ddv7k8ml5gut2.amplifyapp.com/nft-news">Kyzzen</a></b>)
{{range $index, $item := .}}
{{add $index 1}}. {{$item.Title}}
<a href="{{$item.Link}}">Read More - {{capitalize $item.Source}}</a>
{{end}}
			`
	message := helper.FormatHTMLMessage(res, tmpl)

	return message
}

func GetSpaces() string {
	reqData := helper.APIRequest{
		Method: "GET",
		Route:  "fetch-spaces",
	}

	res, err := helper.FetchData(&reqData)
	if err != nil {
		log.Fatalln("Error Occured")
	}

	// Define the HTML template
	const tmpl = `
<b>Upcoming X Spaces Today (from <a href="https://pr-1540.ddv7k8ml5gut2.amplifyapp.com/twitter-spaces">Kyzzen</a></b>)
{{range $index, $item := .}}
{{formatDate $item.Scheduled}} UTC
<b>{{cleanText $item.Title}}</b> <a href="{{$item.Space}}">(View Space)</a>
Host: <a href="x.com/{{$item.UserHandle}}">{{$item.UserHandle}}</a>
{{end}}
<i>Check out the full list of upcoming X spaces on <a href="https://pr-1540.ddv7k8ml5gut2.amplifyapp.com/twitter-spaces">Kyzzen</a>:</i>
			`

	message := helper.FormatHTMLMessage(res[:20], tmpl)

	return message
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
