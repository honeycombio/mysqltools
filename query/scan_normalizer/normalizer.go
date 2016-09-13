package scan_normalizer

import (
	"unicode"
)

type Normalizer struct {
}

func (n *Normalizer) NormalizeQuery(q string) string {

	// three bools to manage our state, in order of priority.
	var escaped bool
	var quoted bool
	var number bool

	var needSpace bool
	var canStartNumber bool
	var quoteRune rune

	var rv []rune

	maybeAddSpace := func() {
		if needSpace {
			needSpace = false
			rv = append(rv, ' ')
		}
	}

	qrunes := []rune(q)

	for _, r := range qrunes {
		if escaped {
			// just skip the next character.  this
			// shouldn't happen in non-quoted contexts.
			escaped = false
			continue
		}
		if quoted {
			if r == quoteRune {
				quoted = false
				canStartNumber = true
			} else if r == '\\' {
				escaped = true
			}
			continue
		}
		if number {
			if (r < '0' || r > '9') && r != '.' {
				number = false
				rv = append(rv, r)
			}
			continue
		}

		// not escaped, quoted, or number.

		if r == '"' || r == '\'' {
			maybeAddSpace()
			quoted = true
			quoteRune = r
			rv = append(rv, '?')
			continue
		}

		// not strictly sufficient - `WHERE c = .4` should be `WHERE c = ?`
		if r >= '0' && r <= '9' {
			maybeAddSpace()
			if canStartNumber {
				number = true
				rv = append(rv, '?')
				continue
			}
		}

		if r == ' ' {
			needSpace = true
			canStartNumber = true
			continue
		} else {
			maybeAddSpace()

			canStartNumber = r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r)

			rv = append(rv, unicode.ToLower(r))

		}
	}

	return string(rv)
}
