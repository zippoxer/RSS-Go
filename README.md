# RSS GO

A simple RSS and ATOM feed reader for the Go programming language.

## Dependencies

- golang.org/x/text/encoding/charmap

## Usage

Send an io.Reader to rss.Get() and it'll do the rest.

### Date/Time Parsing

RSS Go will attempt to parse the date/time values it discovers.

### Encoding

RSS Go currently supports the following encoding:
- CodePage437
- CodePage850
- CodePage852
- CodePage855
- CodePage858
- CodePage862
- CodePage866
- ISO8859-1
- ISO8859-2
- ISO8859-3
- ISO8859-4
- ISO8859-5
- ISO8859-6
- ISO8859-6E
- ISO8859-6I
- ISO8859-7
- ISO8859-8
- ISO8859-8E
- ISO8859-8I
- ISO8859-10
- ISO8859-13
- ISO8859-14
- ISO8859-15
- ISO8859-16
- KOI8R
- KOI8U
- Macintosh
- MacintoshCyrillic
- Windows874
- Windows1250
- Windows1251
- Windows1252
- Windows1253
- Windows1254
- Windows1255
- Windows1256
- Windows1257
- Windows1258