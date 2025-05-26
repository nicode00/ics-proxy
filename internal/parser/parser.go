package parser

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

var filter = Filter{"DTSTART", "greater_than", "20250426"}
var filter2 = Filter{"DTSTART", "less_than", "20250626"}
var filter3 = Filter{"SUMMARY", "not_equals", "OOO"}

type event struct {
	content []string
	fields  map[string]string
}

func Trim(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for {
		ok := scanner.Scan()
		if !ok {
			if scanner.Err() == nil {
				break
			}
			panic("unexpected error while reading ical: " + scanner.Err().Error())
		}
		line := scanner.Text()

		if strings.HasPrefix(line, "BEGIN:VEVENT") {
			e := parseEvent(scanner)
			if filter.shouldInclude(e.fields) && filter2.shouldInclude(e.fields) && filter3.shouldInclude(e.fields) {
				for _, l := range e.content {
					for len(l) > 1 {
						i := min(75, len(l))
						_, err := w.Write([]byte(l[:i] + "\r\n"))
						if err != nil {
							panic("unexpected error while reading ical: " + err.Error())
						}
						l = " " + l[i:]
					}
				}
			}
		} else {
			_, err := w.Write([]byte(line + "\r\n"))
			if err != nil {
				panic("unexpected error while reading ical: " + err.Error())
			}
		}
	}
}

func parseEvent(scanner *bufio.Scanner) event {
	content := []string{"BEGIN:VEVENT"}

	for {
		ok := scanner.Scan()
		if !ok {
			panic("unexpected error while reading event: " + scanner.Err().Error())
		}
		line := scanner.Text()

		if strings.HasPrefix(line, "END:VEVENT") {
			content = append(content, "END:VEVENT")
			break
		}

		if unicode.IsSpace(rune(line[0])) {
			content[len(content)-1] = content[len(content)-1] + strings.TrimSpace(line)
		} else {
			content = append(content, strings.TrimSpace(line))
		}
	}

	fields := make(map[string]string)
	for _, line := range content {
		var field, val string
		for i, ch := range line {
			if ch == ':' || ch == ';' {
				field = line[:i]
				break
			}
		}
		for i := len(line) - 1; i >= 0; i-- {
			if line[i] == ':' || line[i] == ';' {
				val = line[i+1:]
				break
			}
		}
		fields[field] = val
	}

	return event{content, fields}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
