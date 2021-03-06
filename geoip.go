package main

import (
	"encoding/json"
	"fmt"
	"github.com/kballard/goirc/irc"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func init() {
	RegisterModule("geoip", func() Module {
		return &GeoipMod{}
	})
}

//{"ip":"76.126.11.172","country_code":"US","country_name":"United States",
// "region_code":"CA","region_name":"California","city":"San Francisco",
// "zipcode":"94110","latitude":37.7484,"longitude":-122.4156,"metro_code":"807","areacode":"415"}
type FreeGeoip struct {
	Ip           string
	Country_code string
	Country_name string
	Region_code  string
	Region_name  string
	City         string
	Zipcode      string
	Latitude     float64
	Longitude    float64
	Metro_code   string
	Areacode     string
}

type GeoipMod struct {
	urlfmt string
}

func (g *GeoipMod) Init(b *Bot, conn irc.SafeConn) error {
	g.urlfmt = "http://freegeoip.net/%s/%s"

	conn.AddHandler("PRIVMSG", func(c *irc.Conn, l irc.Line) {
		args := strings.Split(l.Args[1], " ")
		if args[0] == ".geo" {
			ip := strings.Join(args[1:], "")
			loc := g.geo(ip)

			if l.Args[0][0] == '#' {
				c.Privmsg(l.Args[0], loc)
			} else {
				c.Privmsg(l.Src.String(), loc)
			}
		}
	})

	log.Printf("geoip module initialized with urlfmt %s", g.urlfmt)

	return nil
}

func (g *GeoipMod) Reload() error {
	return nil
}

func (g *GeoipMod) Call(args ...string) error {
	return nil
}

// return human readable form of geoip data
func (g *GeoipMod) geo(ip string) string {
	var geo FreeGeoip
	var body []byte

	url := fmt.Sprintf(g.urlfmt, "json", ip)

	resp, err := http.Get(url)
	if err != nil {
		goto bad
	}

	defer resp.Body.Close()

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		goto bad
	}

	if err = json.Unmarshal(body, &geo); err != nil {
		goto bad
	}

	return fmt.Sprintf("%s: %s - %s - %s : φ%f° λ%f°",
		geo.Ip, geo.Country_name, geo.Region_name, geo.City, geo.Latitude, geo.Longitude)

bad:
	return fmt.Sprintf("geoip error: %s", err)
}
