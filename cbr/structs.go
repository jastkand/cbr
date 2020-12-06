package cbr

import "encoding/xml"

type DailyResponse struct {
	XMLName xml.Name `xml:"ValCurs" json:"-"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valute  []Valute
}

type Valute struct {
	ID       string `xml:"ID,attr"`
	NumCode  string
	CharCode string
	Nominal  int
	Name     string
	Value    string
}
