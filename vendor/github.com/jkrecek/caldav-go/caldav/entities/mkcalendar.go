package entities

import "encoding/xml"

type MKCalendar struct {
	XMLName xml.Name    `xml:"urn:ietf:params:xml:ns:caldav mkcalendar"`
	Set     *SetPropSet `xml:",omitempty"`
}

type SetPropSet struct {
	XMLName xml.Name `xml:"DAV: set"`
	Props   []*Prop  `xml:",omitempty"`
}

// TODO more parameters
func NewCalendarRequest(name string) *MKCalendar {
	return &MKCalendar{
		Set: &SetPropSet{
			Props: []*Prop{
				{DisplayName: name},
			},
		},
	}
}
