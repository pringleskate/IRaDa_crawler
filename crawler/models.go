package crawler

import "encoding/xml"

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel `xml:"channel"`
	Version string `xml:"version,attr"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Title string `xml:"title"`
	Description string `xml:"description"`
	Link string `xml:"link"`
	Image Image `xml:"image"`
	Atom  Atom `xml:"atom"`
	Item []Item `xml:"item"`
}

type Image struct {
	XMLName xml.Name `xml:"image"`
	URL string `xml:"url"`
	Title string `xml:"title"`
	Link string `xml:"link"`
	Width int `xml:"width"`
	Height int `xml:"height"`
}

type Atom struct {
	XMLName xml.Name `xml:"atom"`
	Rel string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
	Href string `xml:"href,attr"`
}

type Item struct {
	XMLName xml.Name `xml:"item"`
	Guid string `xml:"guid"`
	Title string `xml:"title"`
	Link string `xml:"link"`
	Description string `xml:"description"`
	Pubdate string `xml:"pubDate"`
	Enclosure Enclosure `xml:"enclosure"`
	Category string `xml:"category"`
}

type Enclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	URL string `xml:"url,attr"`
	Type string `xml:"type,attr"`
	Length string `xml:"length,attr"`
}
