package logger

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"
)

func init() {
	for f := range fields {
		Fields = append(Fields, f)
	}
	sort.Strings(Fields)
}

// Fields contains a list of all known static log fields in alphabetical order.
var Fields []string

// pattern is a log output format.
type pattern []field

func (p pattern) write(b *bytes.Buffer, e *Event) {
	for _, fn := range p {
		fn(b, e)
	}
	if b.Len() == 0 {
		return
	}
	b.WriteRune('\n')
}

// field renders a part of the log line.
type field func(b *bytes.Buffer, e *Event)

// fields contains the known log fields and their field functions. The field
// functions should avoid to alloc memory at all cost since they are in the hot
// path. Do not use fmt.Sprintf() but combine the value from the parts. Instead
// of strconv.Atoi/FormatInt() use the local atoi() function which does not
// alloc.
var fields = map[string]field{
	"$remote_addr": func(b *bytes.Buffer, e *Event) {
		if e.Request == nil {
			return
		}
		b.WriteString(e.Request.RemoteAddr)
	},
	"$remote_host": func(b *bytes.Buffer, e *Event) {
		if e.Request == nil {
			return
		}
		host, _ := hostport(e.Request.RemoteAddr)
		b.WriteString(host)
	},
	"$remote_port": func(b *bytes.Buffer, e *Event) {
		if e.Request == nil {
			return
		}
		_, port := hostport(e.Request.RemoteAddr)
		b.WriteString(port)
	},
	"$request": func(b *bytes.Buffer, e *Event) {
		if e.Request == nil {
			return
		}
		b.WriteString(e.Request.Method)
		b.WriteRune(' ')
		b.WriteString(e.Request.RequestURI)
		b.WriteRune(' ')
		b.WriteString(e.Request.Proto)
	},
	"$request_args": func(b *bytes.Buffer, e *Event) {
		// cannot use e.Req.URL since it may have been modified
		if e.RequestURL == nil {
			return
		}
		b.WriteString(e.RequestURL.RawQuery)
	},
	"$request_host": func(b *bytes.Buffer, e *Event) {
		if e.Request == nil {
			return
		}
		b.WriteString(e.Request.Host)
	},
	"$request_method": func(b *bytes.Buffer, e *Event) {
		if e.Request == nil {
			return
		}
		b.WriteString(e.Request.Method)
	},
	"$request_scheme": func(b *bytes.Buffer, e *Event) {
		// cannot use e.Req.URL since it may have been modified
		if e.RequestURL == nil {
			return
		}
		b.WriteString(e.RequestURL.Scheme)
	},
	"$request_uri": func(b *bytes.Buffer, e *Event) {
		if e.Request == nil {
			return
		}
		b.WriteString(e.Request.RequestURI)
	},
	"$request_url": func(b *bytes.Buffer, e *Event) {
		// cannot use e.Req.URL since it may have been modified
		if e.RequestURL == nil {
			return
		}
		b.WriteString(e.RequestURL.String())
	},
	"$request_proto": func(b *bytes.Buffer, e *Event) {
		if e.Request == nil {
			return
		}
		b.WriteString(e.Request.Proto)
	},
	"$response_body_size": func(b *bytes.Buffer, e *Event) {
		atoi(b, e.Response.ContentLength, 0)
	},
	"$response_status": func(b *bytes.Buffer, e *Event) {
		atoi(b, int64(e.Response.StatusCode), 0)
	},
	"$response_time_ms": func(b *bytes.Buffer, e *Event) {
		d := e.End.Sub(e.Start).Nanoseconds()
		s, us := d/int64(time.Second), d%int64(time.Second)/int64(time.Millisecond)
		atoi(b, s, 0)
		b.WriteRune('.')
		atoi(b, us, 3)
	},
	"$response_time_us": func(b *bytes.Buffer, e *Event) {
		d := e.End.Sub(e.Start).Nanoseconds()
		s, us := d/int64(time.Second), d%int64(time.Second)/int64(time.Microsecond)
		atoi(b, s, 0)
		b.WriteRune('.')
		atoi(b, us, 6)
	},
	"$response_time_ns": func(b *bytes.Buffer, e *Event) {
		d := e.End.Sub(e.Start).Nanoseconds()
		s, ns := d/int64(time.Second), d%int64(time.Second)/int64(time.Nanosecond)
		atoi(b, s, 0)
		b.WriteRune('.')
		atoi(b, ns, 9)
	},
	"$time_unix_ms": func(b *bytes.Buffer, e *Event) {
		atoi(b, e.End.UnixNano()/int64(time.Millisecond), 0)
	},
	"$time_unix_us": func(b *bytes.Buffer, e *Event) {
		atoi(b, e.End.UnixNano()/int64(time.Microsecond), 0)
	},
	"$time_unix_ns": func(b *bytes.Buffer, e *Event) {
		atoi(b, e.End.UnixNano(), 0)
	},
	"$time_common": func(b *bytes.Buffer, e *Event) {
		atoi(b, int64(e.End.Day()), 2)
		b.WriteRune('/')
		b.WriteString(shortMonthNames[e.End.Month()])
		b.WriteRune('/')
		atoi(b, int64(e.End.Year()), 4)
		b.WriteRune(':')
		atoi(b, int64(e.End.Hour()), 2)
		b.WriteRune(':')
		atoi(b, int64(e.End.Minute()), 2)
		b.WriteRune(':')
		atoi(b, int64(e.End.Second()), 2)
		b.WriteString(" +0000") // TODO(fs): local time
	},
	"$time_rfc3339": func(b *bytes.Buffer, e *Event) {
		atoi(b, int64(e.End.Year()), 4)
		b.WriteRune('-')
		atoi(b, int64(e.End.Month()), 2)
		b.WriteRune('-')
		atoi(b, int64(e.End.Day()), 2)
		b.WriteRune('T')
		atoi(b, int64(e.End.Hour()), 2)
		b.WriteRune(':')
		atoi(b, int64(e.End.Minute()), 2)
		b.WriteRune(':')
		atoi(b, int64(e.End.Second()), 2)
		b.WriteRune('Z')
	},
	"$time_rfc3339_ms": func(b *bytes.Buffer, e *Event) {
		atoi(b, int64(e.End.Year()), 4)
		b.WriteRune('-')
		atoi(b, int64(e.End.Month()), 2)
		b.WriteRune('-')
		atoi(b, int64(e.End.Day()), 2)
		b.WriteRune('T')
		atoi(b, int64(e.End.Hour()), 2)
		b.WriteRune(':')
		atoi(b, int64(e.End.Minute()), 2)
		b.WriteRune(':')
		atoi(b, int64(e.End.Second()), 2)
		b.WriteRune('.')
		atoi(b, int64(e.End.Nanosecond())/int64(time.Millisecond), 3)
		b.WriteRune('Z')
	},
	"$time_rfc3339_us": func(b *bytes.Buffer, e *Event) {
		atoi(b, int64(e.End.Year()), 4)
		b.WriteRune('-')
		atoi(b, int64(e.End.Month()), 2)
		b.WriteRune('-')
		atoi(b, int64(e.End.Day()), 2)
		b.WriteRune('T')
		atoi(b, int64(e.End.Hour()), 2)
		b.WriteRune(':')
		atoi(b, int64(e.End.Minute()), 2)
		b.WriteRune(':')
		atoi(b, int64(e.End.Second()), 2)
		b.WriteRune('.')
		atoi(b, int64(e.End.Nanosecond())/int64(time.Microsecond), 6)
		b.WriteRune('Z')
	},
	"$time_rfc3339_ns": func(b *bytes.Buffer, e *Event) {
		atoi(b, int64(e.End.Year()), 4)
		b.WriteRune('-')
		atoi(b, int64(e.End.Month()), 2)
		b.WriteRune('-')
		atoi(b, int64(e.End.Day()), 2)
		b.WriteRune('T')
		atoi(b, int64(e.End.Hour()), 2)
		b.WriteRune(':')
		atoi(b, int64(e.End.Minute()), 2)
		b.WriteRune(':')
		atoi(b, int64(e.End.Second()), 2)
		b.WriteRune('.')
		atoi(b, int64(e.End.Nanosecond()), 9)
		b.WriteRune('Z')
	},
	"$upstream_addr": func(b *bytes.Buffer, e *Event) {
		b.WriteString(e.UpstreamAddr)
	},
	"$upstream_host": func(b *bytes.Buffer, e *Event) {
		host, _ := hostport(e.UpstreamAddr)
		b.WriteString(host)
	},
	"$upstream_port": func(b *bytes.Buffer, e *Event) {
		_, port := hostport(e.UpstreamAddr)
		b.WriteString(port)
	},
	"$upstream_request_scheme": func(b *bytes.Buffer, e *Event) {
		if e.UpstreamURL == nil {
			return
		}
		b.WriteString(e.UpstreamURL.Scheme)
	},
	"$upstream_request_uri": func(b *bytes.Buffer, e *Event) {
		if e.UpstreamURL == nil {
			return
		}
		b.WriteString(e.UpstreamURL.RequestURI())
	},
	"$upstream_request_url": func(b *bytes.Buffer, e *Event) {
		if e.UpstreamURL == nil {
			return
		}
		b.WriteString(e.UpstreamURL.String())
	},
	"$upstream_service": func(b *bytes.Buffer, e *Event) {
		b.WriteString(e.UpstreamService)
	},
}

var shortMonthNames = []string{
	"---",
	"Jan",
	"Feb",
	"Mar",
	"Apr",
	"May",
	"Jun",
	"Jul",
	"Aug",
	"Sep",
	"Oct",
	"Nov",
	"Dec",
}

// hostport is a simplified no-alloc version of
// net.SplitHostPort. Since we know that the
// address values have the correct form we can
// skip all the error checking.
func hostport(s string) (host, port string) {
	if s == "" {
		return "", ""
	}
	n := strings.LastIndexByte(s, ':')
	return s[:n], s[n+1:]
}

// atoi is a replacement for strconv.Atoi/strconv.FormatInt
// which does not alloc.
func atoi(b *bytes.Buffer, i int64, pad int) {
	var flag bool
	if i < 0 {
		flag = true
		i = -i
	}

	// format number
	// 2^63-1 == 9223372036854775807
	var d [128]byte
	n, p := len(d), len(d)-1
	for i >= 0 {
		d[p] = byte('0') + byte(i%10)
		i /= 10
		p--
		if i == 0 {
			break
		}
	}

	// padding
	for n-p-1 < pad {
		d[p] = byte('0')
		p--
	}

	if flag {
		d[p] = '-'
		p--
	}
	b.Write(d[p+1:])
}

// parse parses a format string into a pattern based on the following rules:
//
// The format string consists of text and fields. Field names start with a '$'
// and consist of ASCII characters [a-zA-Z0-9.-_]. Field names like
// '$header.name' will render the HTTP header 'name'. All other field names
// must exist in the fields map.
func parse(format string, fields map[string]field) (p pattern, err error) {
	// text is a helper to add raw text to the log output.
	text := func(s string) field {
		return func(b *bytes.Buffer, e *Event) {
			b.WriteString(s)
		}
	}

	// header is a helper to add an HTTP header to the log output.
	header := func(name string) field {
		return func(b *bytes.Buffer, e *Event) {
			if e.Request == nil || e.Request.Header == nil {
				return
			}
			b.WriteString(e.Request.Header.Get(name))
		}
	}

	s := []rune(format)
	for {
		if len(s) == 0 {
			break
		}
		typ, n := lex(s)
		val := string(s[:n])
		s = s[n:]
		switch typ {
		case itemText:
			p = append(p, text(val))
		case itemHeader:
			p = append(p, header(val[len("$header."):]))
		case itemField:
			f := fields[val]
			if f == nil {
				return nil, fmt.Errorf("invalid field %q", val)
			}
			p = append(p, f)
		}
	}
	return p, nil
}

type itemType int

const (
	itemText itemType = iota
	itemField
	itemHeader
)

func (t itemType) String() string {
	switch t {
	case itemText:
		return "TEXT"
	case itemField:
		return "FIELD"
	case itemHeader:
		return "HEADER"
	}
	panic("invalid")
}

type state int

const (
	stateStart state = iota
	stateText
	stateDollar
	stateField
	stateDot
	stateHeader
)

func lex(s []rune) (typ itemType, n int) {
	isIDChar := func(r rune) bool {
		return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9' || r == '_' || r == '-'
	}

	state := stateStart
	for i, r := range s {
		switch state {
		case stateStart:
			switch r {
			case '$':
				state = stateDollar
			default:
				state = stateText
			}

		case stateText:
			switch r {
			case '$':
				return itemText, i
			default:
				// state = stateText
			}

		case stateDollar:
			switch {
			case isIDChar(r):
				state = stateField
			default:
				state = stateText
			}

		case stateField:
			switch {
			case r == '.':
				if string(s[:i]) == "$header" {
					state = stateDot
				} else {
					return itemField, i
				}
			case isIDChar(r):
				// state = stateField
			default:
				return itemField, i
			}

		case stateDot:
			switch {
			case isIDChar(r):
				state = stateHeader
			default:
				return itemField, i
			}

		case stateHeader:
			switch {
			case isIDChar(r):
				// state = stateHeader
			default:
				return itemHeader, i
			}
		}
	}

	switch state {
	case stateDot:
		return itemField, len(s) - 1
	case stateField:
		return itemField, len(s)
	case stateHeader:
		return itemHeader, len(s)
	default:
		return itemText, len(s)
	}
}
