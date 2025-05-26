package matcher

type Matcher struct {
	target []byte
	next   int
}

func New(s string) Matcher {
	return Matcher{target: []byte(s)}
}

func (m *Matcher) Match(ch byte) bool {
	if m.target[m.next] == ch {
		m.next += 1
	} else {
		m.next = 0
	}

	if m.next == len(m.target) {
		return true
	}

	return false
}
