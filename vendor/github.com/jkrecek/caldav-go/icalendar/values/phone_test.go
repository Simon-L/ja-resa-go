package values

import (
	"github.com/jkrecek/caldav-go/icalendar"
	. "gopkg.in/check.v1"
	"testing"
)

type phoneSuite struct{}

type phoneTestObj struct {
	*Phone
}

var _ = Suite(new(phoneSuite))

func TestPhone(t *testing.T) { TestingT(t) }

func (s *phoneSuite) TestEncode(c *C) {
	eto := new(phoneTestObj)
	eto.Phone = &Phone{
		Number: "111 222 333",
		Types: []string{
			"CELL",
			"VOICE",
		},
		IsPreferred: true,
	}

	encoded, err := icalendar.Marshal(eto)
	c.Assert(err, IsNil)
	expected := "TEL;TYPE=CELL;TYPE=VOICE;TYPE=pref:111 222 333"
	c.Assert(encoded, Equals, expected)
}
