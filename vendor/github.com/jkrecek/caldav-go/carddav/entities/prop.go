package entities

import "encoding/xml"

type Prop struct {
	XMLName xml.Name `xml:"DAV: prop"`

	GetETag     *GetETag     `xml:",omitempty"`
	AddressData *AddressData `xml:",omitempty"`
}

type GetETag struct {
	XMLName xml.Name `xml:"DAV: getetag"`
}

type CProp struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:carddav prop"`
	Name    string   `xml:"name,attr"`
}
