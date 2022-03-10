package sql

import (
	"fmt"
	"net"
	"strings"

	"github.com/miekg/dns"
)

type Record struct {
	ID       uint   `gorm:"primary_key"`
	Name     string `gorm:"type:varchar(255)"`
	Type     string `gorm:"type:varchar(10)"`
	Content  string `gorm:"type:text"`
	Ttl      uint32
	Disabled bool
}

func (r *Record) ToRR(qName string, qClass uint16) (dns.RR, error) {
	typ := dns.StringToType[r.Type]
	hrd := dns.RR_Header{Name: qName, Rrtype: typ, Class: qClass, Ttl: r.Ttl}
	if !strings.HasSuffix(hrd.Name, ".") {
		hrd.Name += "."
	}
	rr := dns.TypeToRR[typ]()

	// TODO: support more type
	// this is enough for most query
	switch rr := rr.(type) {
	case *dns.SOA:
		rr.Hdr = hrd
		if !ParseSOA(rr, r.Content) {
			rr = nil
		}
	case *dns.A:
		rr.Hdr = hrd
		rr.A = net.ParseIP(r.Content)
	case *dns.AAAA:
		rr.Hdr = hrd
		rr.AAAA = net.ParseIP(r.Content)
	case *dns.TXT:
		rr.Hdr = hrd
		rr.Txt = []string{r.Content}
	case *dns.NS:
		rr.Hdr = hrd
		rr.Ns = r.Content
	case *dns.PTR:
		rr.Hdr = hrd
		rr.Ptr = r.Content
	default:
		return nil, fmt.Errorf("Unsuported record: %+v", r)
	}

	return rr, nil
}
