package conf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	DefaultSource = &DummySource{}

	assert.Equal(t, "", Get("Abc"))
	assert.Equal(t, false, GetBool("Abc", false))
}
