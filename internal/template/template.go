package template

const News = `
<b>Latest NFT News (from <a href="https://kyzzen.io/nft-news">Kyzzen</a></b>)
{{range $index, $item := .}}
{{add $index 1}}. {{$item.Title}}
<a href="{{$item.Link}}">Read More - {{capitalize $item.Source}}</a>
{{end}}
			`

	// Define the HTML template
	const Spaces = `
<b>Upcoming X Spaces Today (from <a href="https://kyzzen.io/twitter-spaces">Kyzzen</a></b>)
{{range $index, $item := .}}
{{formatDate $item.Scheduled}} UTC
<b>{{cleanText $item.Title}}</b> <a href="{{$item.Space}}">(View Space)</a>
Host: <a href="x.com/{{$item.UserHandle}}">{{$item.UserHandle}}</a>
{{end}}
<i>Check out the full list of upcoming X spaces on <a href="https://kyzzen.io/twitter-spaces">Kyzzen</a>:</i>
			`