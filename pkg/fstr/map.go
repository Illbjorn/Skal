package fstr

import "bytes"

/*
Map allows input of a format string incorporating a string interpolation
syntax where the interpolation signaller is a key enclosed in curly braces
(example: "Hello, {name}."). Values to interpolate to the input string may then
be provided via map arg `vs`.

Example:

	m := map[string]string{
		"name1": "John",
		"name2": "Bill",
	}

	v := Map("Hello, {name1} and {name2}!", m)

	fmt.Println(v) // Hello, John and Bill!
*/
//goland:noinspection GoUnusedExportedFunction
func Map(v string, vs map[string]string) string {
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

		// If we found an opening brace and it's not escaped, consume it as an
		// interpolation token.
		if r == '{' && (pos-1 < 0 || runes[pos-1] != '\\') {
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

			// Look for a hit in our interpolation map.
			if v, ok := vs[token]; ok {
				// If we find one, continue on to avoid the raw token getting written
				// back out.
				out.WriteString(v)
				continue outer
			}

			out.WriteString(token)
		} else {
			out.WriteRune(runes[pos+1])
			pos++
		}
	}

	return out.String()
}
