package fstr

import (
	"bytes"
	"strings"
)

/*
Pairs allows input of a format string incorporating a string interpolation
syntax where the interpolation signaller is a key enclosed in curly braces
(example: "Hello, {name}."). Values to interpolate to the input string may then
be provided to variadic arg `vs` in groups of two where for each two inputs the
first represents a key and the second represents the correlative value.

Example:

	v := Pairs(
		"Hello, {name1} and {name2}!",
		"name1", "John",
		"name2", "Bill",
	)
	fmt.Println(v) // Hello, John and Bill!
*/
func Pairs(v string, vs ...string) string {
	// Create a rune slice from the input string.
	runes := []rune(v)

	// Initialize the buffer to write the resulting string to.
	out := bytes.NewBuffer(nil)

	// Track input string position.
	pos := -1

	// Loop until we hit the end of the line.
outer:
	for pos < len(runes)-1 {
		// Lookahead one rune.
		r := runes[pos+1]

		if r == '\\' {
			if pos+2 < len(runes) {
				if runes[pos+2] == '{' {
					pos++
					pos++
					out.WriteRune(runes[pos])
					continue
				}
			}
		}

		// If we found an opening brace and it's not escaped, consume it as an
		// interpolation token.
		if r == '{' {
			// {
			pos++

			// Consume the interpolation token.
			start := pos + 1
			for runes[pos+1] != '}' {
				pos++
			}
			token := string(runes[start : pos+1])

			// }
			pos++

			// Look for a hit in our interpolation pairs.
			for i := 0; i < len(vs); i += 2 {
				if vs[i] == token {
					// If we find one, continue on to avoid the raw token getting written
					// back out.
					out.WriteString(vs[i+1])
					continue outer
				}
			}
			out.WriteString(token)
		} else {
			pos++
			out.WriteRune(runes[pos])
		}
	}

	return out.String()
}

func PairsStrip(v string, vs ...string) string {
	// Create a rune slice from the input string.
	runes := []rune(v)

	// Initialize the buffer to write the resulting string to.
	out := bytes.NewBuffer(nil)

	// Track input string position.
	pos := -1

	// Loop until we hit the end of the line.
outer:
	for pos < len(runes)-1 {
		// Lookahead one rune.
		r := runes[pos+1]

		if r == '\\' {
			if pos+2 < len(runes) {
				if runes[pos+2] == '{' {
					pos++
					pos++
					out.WriteRune(runes[pos])
					continue
				}
			}
		}

		// If we found an opening brace and it's not escaped, consume it as an
		// interpolation token.
		if r == '{' {
			// {
			pos++

			// Consume the interpolation token.
			start := pos + 1
			for runes[pos+1] != '}' {
				pos++
			}
			token := string(runes[start : pos+1])

			// }
			pos++

			// Look for a hit in our interpolation pairs.
			for i := 0; i < len(vs); i += 2 {
				if vs[i] == token {
					// If we find one, continue on to avoid the raw token getting written
					// back out.
					out.WriteString(vs[i+1])
					continue outer
				}
			}
			out.WriteString(token)
		} else {
			pos++
			out.WriteRune(runes[pos])
		}
	}

	return strings.Trim(out.String(), "\n")
}
