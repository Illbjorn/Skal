package token

import (
	"fmt"
	"strings"
	"sync"
)

func NewCollection(file, in string) *Collection {
	return &Collection{
		file:   file,
		src:    in,
		tokens: make([]Token, 0),
		// We always look one token ahead for all parsing considerations,
		// so we start at -1 since the first token we'll check is n+1 or index 0.
		pos: -1,
		wg:  &sync.WaitGroup{},
	}
}

type Collection struct {
	file        string
	wg          *sync.WaitGroup
	src         string
	lineMarkers [][3]int
	tokens      []Token
	pos         int
}

func (tc *Collection) MarkLine(lineNum, posStart, posEnd int) {
	tc.lineMarkers = append(tc.lineMarkers, [3]int{lineNum, posStart, posEnd})
}

// Retrieves the entire source text line by a provided line number.
func (tc *Collection) SrcLine(line int) string {
	if line < 0 || line >= len(tc.lineMarkers) {
		return ""
	}

	// Locate the line markers for the current line.
	var start, end int
	for _, position := range tc.lineMarkers {
		if position[0] == line {
			start, end = position[1], position[2]
			break
		}
	}

	// Extract the source text line.
	src := tc.src[start:end]

	// Strip any newlines.
	src = strings.Trim(src, "\n")

	return src
}

// Wait simply blocks on the internal WaitGroup.
func (tc *Collection) Wait() {
	tc.wg.Wait()
}

// InputStream returns a Token channel it then spawns a Goroutine to
// continually consume from. All received Collection are added to the
// Collection's tokens slice.
func (tc *Collection) InputStream() chan Token {
	ch := make(chan Token, 1)
	tc.wg.Add(1)

	// Start the token collector.
	go func() {
		defer tc.wg.Done()
		for tk := range ch {
			if tk == nil {
				continue
			}

			tc.tokens = append(tc.tokens, tk)
		}
	}()

	// Return the input channel.
	return ch
}

// Cur returns the current Token.
func (tc *Collection) Cur() Token {
	// Watch out for the end of the slice.
	if tc.pos >= len(tc.tokens) || tc.pos < 0 {
		return &token{_type: EOF}
	}

	return tc.tokens[tc.pos]
}

// Adv advances the position index ahead 1 and returns a pointer to the new
// "current" Token in the slice.
func (tc *Collection) Adv() Token {
	// Advance forward by 1.
	tc.pos++

	// Watch out for the end of the slice.
	if tc.pos >= len(tc.tokens) {
		return &token{_type: EOF}
	}

	if tc.tokens[tc.pos].Value() == "return" {
		fmt.Println("Boooop")
		fmt.Println("Found!")
		fmt.Println(tc.file)
		fmt.Println(tc.pos)
		fmt.Println("Type:", tc.tokens[tc.pos].Type().String())
	}

	// Return the new "current" token.
	return tc.tokens[tc.pos]
}

// AdvIf only consumes the next character if it matches a provided `tts` value.
//
// The returned boolean Value indicates whether or not we advanced (e.g. the
// provided TokenType was found and consumed).
func (tc *Collection) AdvIf(tts ...Type) (Token, bool) {
	for i := 0; i < len(tts); i++ {
		if tts[i] == tc.LA().Type() {
			tk := tc.Adv()
			return tk, true
		}
	}
	return NewToken(tc), false
}

// AdvT advances forward expecting a single provided Token type. If the Token
// encountered is not as specified, a panic occurs.
func (tc *Collection) AdvT(tt Type) Token {
	// Get the next token.
	tk := tc.Adv()

	// Confirm we have a match.
	if tk.Type() != tt {
		assertError(
			tk.SrcLine(),
			tk,
			sprintf("Expected %s, found %s.", tt, tk.Type()),
		)
	}

	return tk
}

// AdvOneOfT iterates provided TokenTypes, returning the first matching Token
// found.
//
// If an expected type is not found, an error is thrown.
func (tc *Collection) AdvOneOfT(tts ...Type) Token {
	// Get the next token.
	tk := tc.Adv()

	// Iterate provided TokenTypes.
	for _, expectedTT := range tts {
		// If we find a match, return it.
		if tk.Type() == expectedTT {
			return tk
		}
	}

	//
	// If we made it here, we're in error.
	//

	// Prepare the 'expected' string.
	var expecteds []string
	for _, v := range tts {
		expecteds = append(expecteds, v.String())
	}
	expected := join(expecteds, ", ")

	assertError(
		tk.SrcLine(),
		tk,
		sprintf("Expected %s, found %s", expected, tk.Type()),
	)

	return NewToken(tc)
}

// LA returns the lookahead (pos+1) Token.
func (tc *Collection) LA() Token {
	// Confirm we're not attempting to index outside the slice.
	if tc.pos+1 >= len(tc.tokens) {
		return &token{_type: EOF}
	}

	// Return a pointer to the pos+1 token.
	return tc.tokens[tc.pos+1]
}

// NTT returns a boolean Value indicating if the next Token has the provided
// TokenType.
func (tc *Collection) NTT(tts ...Type) bool {
	for _, tt := range tts {
		if tc.LA().Type() == tt {
			return true
		}
	}
	return false
}

func (tc *Collection) LineAheadContains(tt Type) bool {
	line := tc.LA().LineStart()
	for i := tc.pos + 1; i < len(tc.tokens); i++ {
		if tc.tokens[i].LineStart() != line {
			break
		}

		if tc.tokens[i].Type() == tt {
			return true
		}
	}

	return false
}

// Parses a valid ref, then returns the following, if any, token.
func (tc *Collection) LookPastRef() Token {
	// Snapshot current position.
	i := tc.pos

	// Parse a ref.
	for {
		// ID | 'this'
		if !tc.NTT(ID, This) {
			break
		}
		tc.pos++

		// Look for indexing.
		if tc.NTT(BrackOpen) {
			// [
			tc.pos++

			// Consume the index ref.
			for {
				if tc.LA().Type() == BrackClose {
					break
				}
				tc.pos++

				// Get the next token.
				if !tc.NTT(
					This,
					IntL,
					StrL,
					BoolL,
					ID,
					Dot) {
					break
				}
				tc.pos++
			}

			// ]
			if !tc.NTT(BrackClose) {
				break
			}
			tc.pos++
		}

		// Break once we run out of dot operators.
		if !tc.NTT(Dot) {
			break
		}
		tc.pos++
	}

	// Get the next token.
	next := tc.LA()

	// Revert the snapshot.
	tc.pos = i

	return next
}
