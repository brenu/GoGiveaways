package twitterwrapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTwitterWrapper(t *testing.T) {
	twitterWrapper := NewTwitterWrapper()

	assert.NotNil(t, twitterWrapper.TwitterClient)
	assert.NotNil(t, twitterWrapper.HTTPClient)
}

func TestHandleImagePost(t *testing.T) {
	wrapper := NewTwitterWrapper()

	correctPost := wrapper.HandleImagePost("https://www.gamerpower.com/offers/1/5ec53f97228a6.jpg")

	incorrectPost := wrapper.HandleImagePost("")

	expectedIncorrectResponse := int64(0)

	assert.Equal(t, incorrectPost, expectedIncorrectResponse)

	assert.NotEqual(t, correctPost, expectedIncorrectResponse)
}
