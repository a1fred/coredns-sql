package sql

import (
	"log"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/jinzhu/gorm"
)

func init() {
	caddy.RegisterPlugin("sql", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	backend := SQLBackend{}
	c.Next()
	if !c.NextArg() {
		return plugin.Error("sql", c.ArgErr())
	}
	dialect := c.Val()

	if !c.NextArg() {
		return plugin.Error("sql", c.ArgErr())
	}
	arg := c.Val()

	db, err := gorm.Open(dialect, arg)
	if err != nil {
		return err
	}
	backend.db = db

	for c.NextBlock() {
		x := c.Val()
		switch x {
		case "fallthrough":
			backend.Fall.SetZonesFromArgs(c.RemainingArgs())
		case "debug":
			args := c.RemainingArgs()
			for _, v := range args {
				switch v {
				case "db":
					backend.db = backend.db.Debug()
				}
			}
			backend.Debug = true
			log.Println(Name, "enable log", args)
		case "auto-migrate":
			// currently only use records table
			if err := backend.db.AutoMigrate(&Record{}).Error; err != nil {
				return err
			}
		default:
			return plugin.Error("sql", c.Errf("unexpected '%v' command", x))
		}
	}

	if c.NextArg() {
		return plugin.Error("sql", c.ArgErr())
	}

	dnsserver.
		GetConfig(c).
		AddPlugin(func(next plugin.Handler) plugin.Handler {
			backend.Next = next
			return backend
		})

	return nil
}
