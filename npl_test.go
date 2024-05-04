package nlp

import (
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/require"
)

var tokenizeCases = []struct { // anonymous struct
	text   string
	tokens []string
}{
	{"Who's on first?", []string{"who", "on", "first"}},
	{"", nil},
}

func TestTokenizeTable(t *testing.T) {
	for _, tc := range tokenizeCases {
		t.Run(tc.text, func(t *testing.T) {
			tokens := Tokenize(tc.text)
			require.Equal(t, tc.tokens, tokens)
		})
	}
}

// Exercise: Read test cases from tokenize_cases.toml

type tokenizeCase struct {
	Text   string
	Tokens []string
}

type file struct {
	Cases []tokenizeCase
}

func TestTokenizeTOML(t *testing.T) {
	var cases file
	_, err := toml.DecodeFile("./tokenize_cases.toml", &cases)
	require.NoError(t, err)

	for _, tc := range cases.Cases {
		t.Run(tc.Text, func(t *testing.T) {
			tokens := Tokenize(tc.Text)
			require.Equal(t, tc.Tokens, tokens)
		})
	}
}

func TestTokenize(t *testing.T) {
	text := "What's on second?"
	expected := []string{"what", "on", "second"}
	tokens := Tokenize(text)
	require.Equal(t, expected, tokens)
	/* Before testify
	// if tokens != expected { // Can't compare slices with == in Go (only to nil)
	if !reflect.DeepEqual(expected, tokens) {
		t.Fatalf("expected %#v, got %#v", expected, tokens)
	}
	*/
}

func FuzzTokenize(f *testing.F) {
	f.Fuzz(func(t *testing.T, text string) {
		tokens := Tokenize(text)
		lText := strings.ToLower(text)
		for _, tok := range tokens {
			if !strings.Contains(lText, tok) {
				t.Fatal(tok)
			}
		}
	})
}
