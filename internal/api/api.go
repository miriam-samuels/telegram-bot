package api

import (
	"fmt"
	"time"

	"github.com/miriam-samuels/telegram-bot/internal/helper"
	types "github.com/miriam-samuels/telegram-bot/internal/repository"
	"github.com/miriam-samuels/telegram-bot/internal/template"
)

var Collections []map[string]interface{} // store collections ids

// fetch NFT News and format to html
func GetNftNews() (string, error) {
	reqData := helper.GraphQLRequest{
		Query: `
		query MyQuery($orderBy: String, $limit: Int) { 
			nftNews(orderBy: $orderBy, limit: $limit) {
			  nodes {
				preview
				source
				sourceTitle
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

	message, err := helper.FormatHTMLMessage(res["nftNews"].Nodes, template.News)
	if err != nil {
		return "", fmt.Errorf("error formatting message; %v", err)
	}
	return message, nil
}

// fetch spaces and format to html
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
				state
			  }
			}
		  }`,
		Variables: map[string]string{
			"limit":     "10",
			"orderBy":   "scheduled",
			"scheduled": time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		},
	}

	res, err := helper.FetchGraphQlData(&reqData)

	if err != nil {
		return "", fmt.Errorf("error pulling space data; %v", err)
	}

	message, err := helper.FormatHTMLMessage(res["twitterSpace"].Nodes, template.Spaces)
	if err != nil {
		return "", fmt.Errorf("error formatting message; %v", err)
	}

	return message, nil
}

// fetch spaces and format to html
func GetRaffles(name string) (string, error) {
	reqData := helper.GraphQLRequest{
		Query: `
		query MyQuery($orderBy: String,$collectionName: String $limit: Int) {
			raffle(orderBy: $orderBy, collectionName: $collectionName, limit: $limit) {
			  nodes {
				supply
				startDate
				source
				sold
				prize
				price
				name
				moonrankRank
				link
				id
				howRareRank
				floorPrice
				endDate
				collectionName
			  }
			  totalCount
			}
		  }`,
		Variables: map[string]string{
			"limit":          "15",
			"orderBy":        "-endDate",
			"collectionName": name,
		},
	}

	res, err := helper.FetchGraphQlData(&reqData)

	if err != nil {
		return "", fmt.Errorf("error pulling space data; %v", err)
	}

	var raffles []map[string]interface{}
	for i := len(res["raffle"].Nodes) - 1; i > 0; i-- {
		raffle := res["raffle"].Nodes[i]
		if helper.TimeDiff(raffle["endDate"].(string)) > 0 {
			raffles = append(raffles, raffle)
		}
	}

	if len(raffles) == 0 {
		return "No Active Raffle", nil
	}

	message, err := helper.FormatHTMLMessage(raffles, template.CollectionRaffles)
	if err != nil {
		return "", fmt.Errorf("error formatting message; %v", err)
	}

	return message, nil
}

// fetch loan and format to html
func GetLoansData(name string) (string, error) {
	reqData := helper.GraphQLRequest{
		Query: `
		query LendingPool($limit: Int, $collectionName: String){
			lendingPools(limit:$limit, collectionName: $collectionName){
			   nodes {
				  availableLiquidity
				  collectionKey
				  collectionName
				  depositYieldApy
				  duration
				  highestOffer
				  id
				  lastLoan
				  interestRate
				  marketplace
				  thumbnailUrl
				  totalLiquidity
				  lowestOffer
				  minYieldApy
				}
		   }
		 }`,
		Variables: map[string]string{
			"limit":          "4",
			"collectionName": name,
		},
	}

	res, err := helper.FetchGraphQlData(&reqData)
	if err != nil {
		return "", fmt.Errorf("error pulling news data; %v", err)
	}

	message, err := helper.FormatHTMLMessage(res["lendingPools"].Nodes, template.CollectionLoans)
	if err != nil {
		return "", fmt.Errorf("error formatting message; %v", err)
	}
	return message, nil
}

func GetACollection(name string, tmpl string) (string, error) {
	reqData := helper.GraphQLRequest{
		Query: `
		query MyQuery($name: String, $limit: Int, $verifeyed: Boolean) {
			collections(name: $name, limit: $limit, verifeyed: $verifeyed) {
				nodes {
					floorPrice
					floorPriceDelta
					floorPrice24hAgo
					floorPrice30dAgo
					floorPrice7dAgo
					floorPricePast30dDelta
					floorPricePast24hDelta
					floorPricePast7dDelta
					floorPriceUsd1hAgo
					floorPriceUsdPast1hDelta
					floorPriceWithRoyaltiesAndFees
					volumePast24h
					volumePast24hDelta
					volumePast7dDelta
					volumePast7d
					volumePast30dDelta
					volumePast30d
					listings24hDelta
					listings30dDelta
					listings7dDelta
					description
					name
					salesPast7dDelta
					salesPast7d
					salesPast30dDelta
					salesPast30d
					salesPast24hDelta
					salesPast24h
					totalOwners
					items
					sellNowPrice
					id
					owners24hDelta
					owners30dDelta
					owners7dDelta
					twitter
      				website
					listed
				}
			}
		  }`,
		Variables: map[string]string{
			"limit":     "1",
			"name":      name,
			"verifeyed": "true",
		},
	}

	res, err := helper.FetchGraphQlData(&reqData)

	if err != nil {
		return "", fmt.Errorf("error pulling collection data; %v", err)
	}

	if len(res["collections"].Nodes) == 0 {
		return "Collection not found, did you mean", &types.CustomError{Code: 404, Message: name}
	} else {
		//  slot in market cap
		res["collections"].Nodes[0]["marketcap"] = res["collections"].Nodes[0]["items"].(float64) * res["collections"].Nodes[0]["floorPrice"].(float64)
		res["collections"].Nodes[0]["marketcapusd"] = res["collections"].Nodes[0]["items"].(float64) * res["collections"].Nodes[0]["floorPriceUsd1hAgo"].(float64)

	}

	message, err := helper.FormatHTMLMessage(res["collections"].Nodes, tmpl)
	if err != nil {
		return "", fmt.Errorf("error formatting message; %v", err)
	}
	return message, nil
}

func GetAllCollections() ([]map[string]interface{}, error) {
	reqData := helper.GraphQLRequest{
		Query: `
		query MyQuery($orderBy: String, $verifeyed: Boolean) {
			collections(orderBy: $orderBy, verifeyed: $verifeyed) {
				nodes {
					name
					id
				}
			}
		  }`,
		Variables: map[string]string{
			"verifeyed": "true",
		},
	}

	res, err := helper.FetchGraphQlData(&reqData)

	if err != nil {
		return nil, fmt.Errorf("error pulling collection data; %v", err)
	}

	return res["collections"].Nodes, nil
}

func GetCollectionOffers(collection string) (helper.Node, error) {
	reqData := helper.GraphQLRequest{
		Query: `
		query MyQuery($collectionId: String) {
			offers(collectionId: $collectionId) {
			  totalCount
			}
		  }`,
		Variables: map[string]string{
			"collectionId": collection,
		},
	}
	res, err := helper.FetchGraphQlData(&reqData)

	if err != nil {
		return helper.Node{}, fmt.Errorf("error pulling collection data; %v", err)
	}

	return res["offers"], nil
}
