package sql

import (
	"testing"

	"github.com/miekg/dns"
)

func TestParseSOA(t *testing.T) {
	if ParseSOA(new(dns.SOA), "") == true {
		t.Fatalf("incorrect soa parsed")
	}

	if ParseSOA(new(dns.SOA), "ns1.example.com. support.example.com. 2021010141 7200 3600 1209600 3600") == false {
		t.Fatalf("correct soa not parsed")
	}
}
