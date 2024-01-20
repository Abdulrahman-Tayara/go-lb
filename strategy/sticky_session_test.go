package strategy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasSessionToInt(t *testing.T) {
	session := "ABCDEFG"

	num := hashSessionToInt(session)

	assert.Equal(t, 476, num)
}
