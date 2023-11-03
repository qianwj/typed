package regex

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// IgnoreCase
	// Case insensitivity to match upper and lower cases. For an example, see
	// Perform Case-Insensitive Regular Expression Match.
	IgnoreCase = "i"

	// MultilineMatch
	// For patterns that include anchors (i.e. ^ for the start, $ for the end),
	// match at the beginning or end of each line for strings with multiline values.
	// Without this option, these anchors match at beginning or end of the string.
	// For an example, see Multiline Match for Lines Starting with Specified Pattern.
	// If the pattern contains no anchors or if the string value has no newline
	// characters (e.g. \n), the m option has no effect.
	MultilineMatch = "m"

	// Extended
	// "Extended" capability to ignore all white space characters in the $regex
	// pattern unless escaped or included in a character class.
	// Additionally, it ignores characters in-between and including an un-escaped
	// hash/pound (#) character and the next new line, so that you may include comments
	// in complicated patterns. This only applies to data characters; white space
	// characters may never appear within special character sequences in a pattern.
	// The x option does not affect the handling of the VT character (i.e. code 11).
	Extended = "x"

	// AllowDotChar
	// Allows the dot character (i.e. .) to match all characters including newline characters.
	// For an example, see Use the . Dot Character to Match New Line.
	AllowDotChar = "s"
)

type Matcher struct {
	pattern string
	options map[string]bool
}

func NewMatcher(pattern string) *Matcher {
	return &Matcher{
		options: make(map[string]bool),
		pattern: pattern,
	}
}

func (m *Matcher) IgnoreCase() *Matcher {
	m.options[IgnoreCase] = true
	return m
}

func (m *Matcher) MultilineMatch() *Matcher {
	m.options[MultilineMatch] = true
	return m
}

func (m *Matcher) Extended() *Matcher {
	m.options[Extended] = true
	return m
}

func (m *Matcher) AllowDotChar() *Matcher {
	m.options[AllowDotChar] = true
	return m
}

func (m *Matcher) Compile() primitive.Regex {
	options := ""
	for k, add := range m.options {
		if add {
			options += k
		}
	}
	return primitive.Regex{Pattern: m.pattern, Options: options}
}

func (m *Matcher) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Compile())
}

func (m *Matcher) MarshalBSON() ([]byte, error) {
	return bson.Marshal(m.Compile())
}

func (m *Matcher) UnmarshalBSON(bytes []byte) error {
	var reg primitive.Regex
	if err := bson.Unmarshal(bytes, &reg); err != nil {
		return err
	}
	m.pattern = reg.Pattern
	opts := make(map[string]bool)
	for _, r := range []rune(reg.Options) {
		opts[string(r)] = true
	}
	m.options = opts
	return nil
}
