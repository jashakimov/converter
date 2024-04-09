package elecard

import (
	"encoding/xml"
)

type CreateTaskRequest struct {
	XMLName    xml.Name `xml:"XMLConfig"`
	Text       string   `xml:",chardata"`
	Dispatcher string   `xml:"dispatcher,attr"`
	SetValue   struct {
		Text      string `xml:",chardata"`
		Name      string `xml:"Name,attr"`
		Parameter struct {
			Text string   `xml:",chardata"`
			P    []string `xml:"p"`
		} `xml:"Parameter"`
	} `xml:"SetValue"`
}

type CreateTaskResponse struct {
	XMLName    xml.Name `xml:"XMLConfig"`
	Text       string   `xml:",chardata"`
	Dispatcher string   `xml:"dispatcher,attr"`
	SetValue   struct {
		Text   string `xml:",chardata"`
		Name   string `xml:"Name,attr"`
		RetVal struct {
			Text        string `xml:",chardata"`
			WatchFolder struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
			} `xml:"WatchFolder"`
		} `xml:"RetVal"`
	} `xml:"SetValue"`
}

type GetStatusRequest struct {
	XMLName    xml.Name `xml:"XMLConfig"`
	Text       string   `xml:",chardata"`
	Dispatcher string   `xml:"dispatcher,attr"`
	GetValue   struct {
		Text      string `xml:",chardata"`
		Name      string `xml:"Name,attr"`
		Parameter struct {
			Text string   `xml:",chardata"`
			P    []string `xml:"p"`
		} `xml:"Parameter"`
	} `xml:"GetValue"`
}

type GetStatusResponse struct {
	XMLName    xml.Name `xml:"XMLConfig"`
	Text       string   `xml:",chardata"`
	Dispatcher string   `xml:"dispatcher,attr"`
	GetValue   struct {
		Text   string `xml:",chardata"`
		Name   string `xml:"Name,attr"`
		RetVal string `xml:"RetVal"`
	} `xml:"GetValue"`
}
