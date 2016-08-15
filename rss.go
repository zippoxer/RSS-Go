// RSS/ATOM feed reader.
package rss

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"time"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

// Defines an RSS/Atom feed
type Feed struct {
	Title    string
	Subtitle string
	Link     string
	Items    []*Item
}

// Defines a feed item
type Item struct {
	Id          string
	Title       string
	Description string
	Link        string
	When        time.Time
}

const feedTitle = "title"

const (
	rssChannel     = "channel"
	rssItem        = "item"
	rssLink        = "link"
	rssPubDate     = "pubdate"
	rssDescription = "description"
	rssId          = "guid"
)

const (
	atomSubtitle = "subtitle"
	atomFeed     = "feed"
	atomEntry    = "entry"
	atomLink     = "link"
	atomLinkHref = "href"
	atomUpdated  = "updated"
	atomSummary  = "summary"
	atomId       = "id"
)

const (
	levelFeed = 1
	levelPost = 2
)

var charsetMap = map[string]encoding.Encoding{
	"codepage437":       charmap.CodePage437,
	"codepage850":       charmap.CodePage850,
	"CodePage852":       charmap.CodePage852,
	"CodePage855":       charmap.CodePage855,
	"CodePage858":       charmap.CodePage858,
	"CodePage862":       charmap.CodePage862,
	"CodePage866":       charmap.CodePage866,
	"iso88591":          charmap.ISO8859_1,
	"iso88592":          charmap.ISO8859_2,
	"iso88593":          charmap.ISO8859_3,
	"iso88594":          charmap.ISO8859_4,
	"iso88595":          charmap.ISO8859_5,
	"iso88596":          charmap.ISO8859_6,
	"iso88596E":         charmap.ISO8859_6E,
	"iso88596I":         charmap.ISO8859_6I,
	"iso88597":          charmap.ISO8859_7,
	"iso88598":          charmap.ISO8859_8,
	"iso88598E":         charmap.ISO8859_8E,
	"iso88598I":         charmap.ISO8859_8I,
	"iso885910":         charmap.ISO8859_10,
	"iso885913":         charmap.ISO8859_13,
	"iso885914":         charmap.ISO8859_14,
	"iso885915":         charmap.ISO8859_15,
	"iso885916":         charmap.ISO8859_16,
	"koi8r":             charmap.KOI8R,
	"koi8u":             charmap.KOI8U,
	"macintosh":         charmap.Macintosh,
	"macintoshcyrillic": charmap.MacintoshCyrillic,
	"windows874":        charmap.Windows874,
	"windows1250":       charmap.Windows1250,
	"windows1251":       charmap.Windows1251,
	"windows1252":       charmap.Windows1252,
	"windows1253":       charmap.Windows1253,
	"windows1254":       charmap.Windows1254,
	"windows1255":       charmap.Windows1255,
	"windows1256":       charmap.Windows1256,
	"windows1257":       charmap.Windows1257,
	"windows1258":       charmap.Windows1258,
}

// createCharsetReader Creates a new io.reader for a certain charset
func createCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	key := strings.Replace(strings.Trim(strings.ToLower(charset), " "), "-", "", -1)

	if enc, exists := charsetMap[key]; exists {
		return enc.NewDecoder().Reader(input), nil
	}

	return nil, fmt.Errorf("Unknown charset: %s", charset)
}

// Get parses a RSS/Atom feed
func Get(r io.Reader) (*Feed, error) {
	var tag string
	var atom bool
	var level int

	feed := &Feed{}
	item := &Item{}

	d := xml.NewDecoder(r)
	d.CharsetReader = createCharsetReader
	d.Strict = false

	for {
		token, err := d.Token()

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			tag = strings.ToLower(t.Name.Local)
			switch {
			case tag == atomFeed:
				atom = true
				level = levelFeed
			case tag == rssChannel:
				atom = false
				level = levelFeed
			case (!atom && tag == rssItem) || (atom && tag == atomEntry):
				level = levelPost
				item = &Item{When: time.Now()}
			case atom && tag == atomLink:
				for _, a := range t.Attr {
					if strings.ToLower(a.Name.Local) == atomLinkHref {
						switch level {
						case levelFeed:
							feed.Link = a.Value
						case levelPost:
							item.Link = a.Value
						}
						break
					}
				}
			}
		case xml.EndElement:
			e := strings.ToLower(t.Name.Local)
			if e == atomEntry || e == rssItem {
				feed.Items = append(feed.Items, item)
			}
		case xml.CharData:
			text := string([]byte(t))
			if strings.TrimSpace(text) == "" {
				continue
			}
			switch level {
			case levelFeed:
				switch {
				case tag == feedTitle:
					feed.Title = text
				case (!atom && tag == rssDescription) || (atom && tag == atomSubtitle):
					feed.Subtitle = text
				case !atom && tag == rssLink:
					feed.Link = text
				}
			case levelPost:
				switch {
				case (!atom && tag == rssId) || (atom && tag == atomId):
					item.Id = text
				case tag == feedTitle:
					item.Title = text
				case (!atom && tag == rssDescription) || (atom && tag == atomSummary):
					item.Description = text
				case !atom && tag == rssLink:
					item.Link = text
				case atom && tag == atomUpdated:
					var f string
					switch {
					case strings.HasSuffix(strings.ToUpper(text), "Z"):
						f = "2006-01-02T15:04:05Z"
					default:
						f = "2006-01-02T15:04:05-07:00"
					}
					t, err := time.Parse(f, text)
					if err != nil {
						return nil, err
					}
					item.When = t
				case !atom && tag == rssPubDate:
					var f string
					if strings.HasSuffix(strings.ToUpper(text), "T") {
						f = "Mon, 2 Jan 2006 15:04:05 MST"
					} else {
						f = "Mon, 2 Jan 2006 15:04:05 -0700"
					}
					t, err := time.Parse(f, text)
					if err != nil {
						return nil, err
					}
					item.When = t
				}
			}
		}
	}

	return feed, nil
}
