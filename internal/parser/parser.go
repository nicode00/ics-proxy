package parser

import (
	"io"
	"unicode"
)

const BUF_SIZE = 77

type Parser struct {
	r       io.Reader
	w       io.Writer
	atEvent bool
	event   [][]byte
	filters []Filter
}

var filter = Filter{"DTSTART", "greater_than", "20250426"}
var filter2 = Filter{"DTSTART", "less_than", "20250626"}
var filter3 = Filter{"SUMMARY", "not_equals", "OOO"}

func New(r io.Reader, w io.Writer) Parser {
	return Parser{
		r:       r,
		w:       w,
		event:   make([][]byte, 0),
		filters: []Filter{filter, filter2, filter3}}
}

func (p *Parser) Parse() error {
	buf := make([]byte, BUF_SIZE)

	var i, j, n int
	var err error
	hasMore := true

	for len(buf[:j]) > 0 || hasMore {
		if hasMore {
			n, err = p.r.Read(buf[j:])
			j = j + n
		}
		i = p.parseLine(buf[:j]) % (len(buf) + 1)

		if err == io.EOF {
			hasMore = false
		} else if err != nil {
			return err
		}

		buf = append(buf[i:j], make([]byte, BUF_SIZE)...)
		j -= i
	}

	return nil
}

func (p *Parser) parseLine(buf []byte) int {
	eol := indexOfNewline(buf)

	if string(buf[:eol-1]) == "BEGIN:VEVENT" {
		p.atEvent = true
	}

	if !p.atEvent {
		err := p.writeAll(buf[:eol+1])
		if err != nil {
			panic(err)
		}
	} else if unicode.IsSpace(rune(buf[0])) {
		last := p.event[len(p.event)-1]
		p.event[len(p.event)-1] = append(last[:len(last)-2], buf[1:eol+1]...)
	} else {
		p.event = append(p.event, buf[:eol+1])
	}

	if (string(buf[:eol-1])) == "END:VEVENT" {
		p.processEvent(p.event)
		p.event = make([][]byte, 0)
		p.atEvent = false
	}

	return eol + 1
}

func (p *Parser) processEvent(buf [][]byte) {
	fields := make(map[string]string)
	for _, line := range buf {
		fields[fieldName(line)] = fieldValue(line)
	}

	for _, filter := range p.filters {
		if !filter.FieldsMatch(fields) {
			return
		}
	}

	wbuf := make([]rune, 0, 77*len(buf))
	i := 0
	for _, line := range buf {
		for _, r := range string(line) {
			if r == '\r' || r == '\n' {
				wbuf = append(wbuf, r)
				i = 0
			} else if i == 72 {
				wbuf = append(wbuf, '\r', '\n', ' ', r)
				i = 0
			} else {
				wbuf = append(wbuf, r)
				i += 1
			}
		}
	}

	err := p.writeAll([]byte(string(wbuf)))
	if err != nil {
		panic(err)
	}
}

func (p *Parser) writeAll(buf []byte) error {
	i := 0
	for {
		n, err := p.w.Write(buf[i:])
		if err != nil {
			return err
		}
		if i += n; i >= len(buf) {
			return nil
		}
	}
}

func indexOfNewline(buf []byte) int {
	for i := range buf {
		if buf[i] == '\n' {
			return i
		}
	}
	return -1
}

func fieldName(buf []byte) string {
	i := 0
	for ; ; i += 1 {
		if buf[i] == ';' || buf[i] == ':' {
			break
		}
	}
	return string(buf[:i])
}

func fieldValue(buf []byte) string {
	i := len(buf) - 1
	for ; ; i -= 1 {
		if buf[i] == ';' || buf[i] == ':' {
			break
		}
	}
	return string(buf[i+1:])
}
