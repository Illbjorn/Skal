package lex

const (
	// Used by the lexer to indicate we've reached the end of the File.
	rEOF = '\x00'
)

var alphaRunes = []rune{
	// Lowercase
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	// Uppercase
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	// Only included non-alpha symbol.
	'_',
}

func isAlpha(r rune) bool {
	for _, x := range alphaRunes {
		if x == r {
			return true
		}
	}
	return false
}

var numRunes = []rune{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
}

func isNum(r rune) bool {
	for _, x := range numRunes {
		if x == r {
			return true
		}
	}
	return false
}

func isAlphaNum(r rune) bool {
	return isAlpha(r) || isNum(r)
}

var symbolRunes = []rune{
	// Binary operators
	'+', '-', '/', '*',
	// Comparison Operators
	'=', '<', '>',
	// Misc Characters
	',', '.', '!', '(', ')', '[', ']', '{', '}', '|', '&', ':', ';', '#',
}

func isSymbol(r rune) bool {
	for _, x := range symbolRunes {
		if x == r {
			return true
		}
	}
	return false
}
