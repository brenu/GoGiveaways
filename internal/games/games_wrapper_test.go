package games

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGamesLookUp(t *testing.T) {
	result, err := GamesLookUp()
	assert.Nil(t, err)
	assert.NotNil(t, result)
}
