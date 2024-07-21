package template

// Define the HTML template
const WelcomeMessage = `
<b>Welcome to Kyzzen's Telegram Bot 🫡</b>

Explore the world of Solana NFTs and stay updated in the space:

<b>Search NFT Collections:</b>
/searchc [collection] - collection data & analytics

<b>View Top 5 Solana NFT Collections:</b>
/topvolumesol - by Volume
/toptrendingsol - by Volume 24H %△ 
/topgainerssol - by Floor 24H %△ 
/topmktcapsol - by Market Cap

<b>Stay Updated in the NFT Space:</b>
/news - Catch up on the latest NFT news
/spaces - View upcoming X Spaces on NFTs

<b>Learn more about Kyzzen:</b>
- <a href="https://kyzzen.io">Website</a>
- <a href="https://x.com/Kyzzen_io">Twitter/X</a>
- <a href="https://discord.com/invite/hhbgab4bsD">Discord</a>
`

const ErrorMessage = `
Please wait for awhile before retrying your request.

If you continue to experience difficulties, kindly open a support ticket in our <a href="https://discord.com/invite/hhbgab4bsD">Discord</a>

We apologize for any inconvenience caused 🙏.
`

const News = `
<b>Latest NFT News (from <a href="https://kyzzen.io/nft-news">Kyzzen</a>)</b>
{{range $index, $item := .}}
{{add $index 1}}. {{$item.title}}
<a href="{{$item.link}}">Read More - {{$item.sourceTitle}}</a>
{{end}}
<i>⚠️ As these articles come from 3rd-party sources, please verify that the link looks legitimate before clicking through.</i> 

<i>Check out the full list of NFT News on <a href="https://kyzzen.io/nft-news">Kyzzen</a>:</i>
`

const Spaces = `
<b>Upcoming X Spaces Today (from <a href="https://kyzzen.io/twitter-spaces">Kyzzen</a>)</b>
{{range $index, $item := .}}
{{if isLive $item.scheduled}}
<b>Live Now 🟢</b>
{{else}}
{{formatDate $item.scheduled}} UTC
{{end}}<b>{{cleanText $item.title}}</b> <a href="{{$item.spaceUrl}}">(View Space)</a>
Host: <a href="x.com/{{$item.userhandle}}">{{$item.userhandle}}</a>{{end}}

<i>Check out the full list of upcoming X spaces on <a href="https://kyzzen.io/twitter-spaces">Kyzzen</a>:</i>
`

const CollectionInfo = `
{{range $index, $item := .}}
<b>{{$item.name}} (⛓Solana) </b>
<a href="{{$item.website}}">🌐</a> <a href="{{$item.twitter}}">🐦</a>

{{$item.description}}

Floor (SOL):    {{divide $item.floorPrice 1000000000}} SOL 
Floor (USD):    ${{divide $item.floorPriceUsd1hAgo 1}} 
Supply:            {{addComma $item.items}}
Market Cap:    {{divide $item.marketcap 1000000000}} SOL $({{divide $item.marketcapusd 1}}) 

24H Vol:           {{divide $item.volumePast24h 1000000000}} SOL ({{if greater $item.volumePast24hDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.volumePast24hDelta 100}}%)
24H Sales:       {{divide $item.salesPast24h 1000000000}} SOL ({{if greater $item.salesPast24hDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.salesPast24hDelta 100}}%)

7D Vol:             {{divide $item.volumePast7d 1000000000}} SOL ({{if greater $item.volumePast7dDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.volumePast7dDelta 100}}%)
7D Sales:         {{divide $item.salesPast7d 1000000000}} SOL ({{if greater $item.salesPast7dDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.salesPast7dDelta 100}}%)

30D Vol:           {{divide $item.volumePast30d 1000000000}} SOL ({{if greater $item.volumePast30dDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.volumePast30dDelta 100}}%)
30D Sales:       {{divide $item.salesPast30d 1000000000}} SOL ({{if greater $item.salesPast30dDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.salesPast30dDelta 100}}%)

Listings:           {{addComma $item.listed}}
Holders:           {{addComma $item.totalOwners}}
{{end}}
`

const CollectionFloor = `
{{range $index, $item := .}}
<b>Current Floor Price: </b> {{divide $item.floorPrice 1000000000}} SOL / {{divide $item.floorPriceUsd1hAgo 1}} USD 
	• 24H ago:   {{divide $item.floorPrice24hAgo 1000000000}} SOL ({{if greater $item.floorPricePast24hDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.floorPricePast24hDelta 100}}%)
	• 7D ago:     {{divide $item.floorPrice7dAgo 1000000000}} SOL ({{if greater $item.floorPricePast7dDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.floorPricePast7dDelta 100}}%)
	• 30D ago:   {{divide $item.floorPrice30dAgo 1000000000}} SOL ({{if greater $item.floorPricePast30dDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.floorPricePast30dDelta 100}}%)

<b>Sell Now:</b>   {{divide $item.sellNowPrice 1000000000}} SOL
{{end}}
`

const CollectionListing = `
{{range $index, $item := .}}
<b>Current Listings: </b> {{addComma $item.listed}}
	• 24H ago:   {{addComma $item.listed}} ({{if greater $item.listings24hDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.listings24hDelta 100}}%)
	• 7D ago:     {{addComma $item.listed}} ({{if greater $item.listings7dDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.listings7dDelta 100}}%)
	• 30D ago:   {{addComma $item.listed}} ({{if greater $item.listings30dDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.listings30dDelta 100}}%)

{{end}}
`

const CollectionVol = `
{{range $index, $item := .}}
<b>Volume: </b>
	• 24H:   {{divide $item.volumePast24h 1000000000}} SOL ({{if greater $item.volumePast24hDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.volumePast24hDelta 100}}%)
	• 7D:     {{divide $item.volumePast7d 1000000000}} SOL ({{if greater $item.volumePast7dDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.volumePast7dDelta 100}}%)
	• 30D:   {{divide $item.volumePast30d 1000000000}} SOL ({{if greater $item.volumePast30dDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.volumePast30dDelta 100}}%)

<b>No. of New Sales:</b>
	• 24H:    {{divide $item.salesPast24h 1000000000}} SOL ({{if greater $item.salesPast24hDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.salesPast24hDelta 100}}%)
	• 7D:      {{divide $item.salesPast7d 1000000000}} SOL ({{if greater $item.salesPast7dDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.salesPast7dDelta 100}}%)
	• 30D:    {{divide $item.salesPast30d 1000000000}} SOL ({{if greater $item.salesPast30dDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.salesPast30dDelta 100}}%)
{{end}}
`
const CollectionHolders = `
{{range $index, $item := .}}
<b>Current Owners: </b> {{$item.totalOwners}}
	• 24H ago:   {{addComma $item.totalOwners}} ({{if greater $item.owners24hDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.owners24hDelta 100}}%)
	• 7D ago:     {{addComma $item.totalOwners}} ({{if greater $item.owners7dDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.owners7dDelta 100}}%)
	• 30D ago:   {{addComma $item.totalOwners}} ({{if greater $item.owners30dDelta 0}}⬆️{{else}}🔻{{end}} {{divide $item.owners30dDelta 100}}%)
{{end}}
`

const CollectionLoans = `
Highest Loan Offers:
{{range $index, $item := .}}
{{$item.marketplace}}:  {{divide $item.highestOffer 1000000000}} SOL | {{divide $item.depositYieldApy 100}}% APY | {{poolDuration $item}}
{{end}}
`

const CollectionRaffles = `
{{range $index, $item := .}}
{{if eq $index 0}}
<b>{{$item.collectionName}}</b>
{{end}}
<b>Raffle: </b> {{$item.name}}
Rarity Rank :      Moonrank {{$item.moonrankRank}} | HowRare {{$item.howRareRank}}
Price/Ticket:     {{divide $item.price 1000000000}} SOL
Tickets Left:      {{minus $item.supply $item.sold}} / {{$item.supply}}
Ending In:          {{timeTo $item.endDate}}
<a href="{{$item.link}}">View Raffle</a> {{end}}
`

// View more NFT loan offers on <a href="https://kyzzen.io/nft-lending">Kyzzen</a>
