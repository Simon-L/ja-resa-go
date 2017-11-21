package entities

import (
	"encoding/xml"
	. "gopkg.in/check.v1"
	"testing"
)

type ContactQuerySuite struct{}

var _ = Suite(new(ContactQuerySuite))

func TestContactQuery(t *testing.T) { TestingT(t) }

func (s *ContactQuerySuite) TestNewEventRangeQuery(c *C) {
	rq := NewContactQueryWithProps("VERSION", "UID")
	bts, err := xml.Marshal(rq)
	c.Assert(err, IsNil)
	c.Assert(string(bts), Equals, `<addressbook-query xmlns="urn:ietf:params:xml:ns:carddav"><prop xmlns="DAV:"><getetag xmlns="DAV:"></getetag><address-data xmlns="urn:ietf:params:xml:ns:carddav"><prop xmlns="urn:ietf:params:xml:ns:carddav" name="VERSION"></prop><prop xmlns="urn:ietf:params:xml:ns:carddav" name="UID"></prop></address-data></prop></addressbook-query>`)
}
