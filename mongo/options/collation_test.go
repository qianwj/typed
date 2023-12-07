package options

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollation_Locale(t *testing.T) {
	expected := "EN"
	coll := NewCollation().Locale(expected)
	actual := coll.Raw().Locale
	assert.Equal(t, expected, actual)
}
