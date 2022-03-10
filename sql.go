// Package sql implements a plugin that query sql database to resolve the coredns query
package sql

import (
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/fall"
	"github.com/coredns/coredns/request"
	"github.com/jinzhu/gorm"
	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

const Name = "sql"

type SQLBackend struct {
	db    *gorm.DB
	Debug bool
	Next  plugin.Handler
	Fall  fall.F
}

func (sqldb SQLBackend) Name() string { return Name }
func (sqldb SQLBackend) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	a := new(dns.Msg)
	a.SetReply(r)
	a.Compress = true
	a.Authoritative = true

	records, err := FindSql(sqldb.db, state.QName(), state.Type())
	if err != nil {
		return dns.RcodeServerFailure, err
	}
	if len(records) > 0 {
		for _, r := range records {
			rr, err := r.ToRR(state.QName(), state.QClass())
			if err != nil {
				return dns.RcodeServerFailure, err
			}

			a.Answer = append(a.Answer, rr)
		}
		return dns.RcodeSuccess, w.WriteMsg(a)
	}

	if sqldb.Fall.Through(state.Name()) {
		return plugin.NextOrFailure(sqldb.Name(), sqldb.Next, ctx, w, r)
	}

	// TODO: maybe reply with soa
	return 0, w.WriteMsg(a)
}
