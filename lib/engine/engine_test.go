package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEngine_New(t *testing.T) {
	_, err := New("echo", "Hello World")
	assert.Nil(t, err)
}

func TestEngine_Wait(t *testing.T) {
	en, err := New("seq", "10")
	assert.Nil(t, err)

	err = en.Wait("8")
	assert.Nil(t, err)
}

func TestEngine_SetOption(t *testing.T) {
	en, err := New("cat")
	assert.Nil(t, err)

	err = en.Set(map[string]string{
		"USI_Hash": "334",
	})
	assert.Nil(t, err)

	err = en.Wait("setoption name USI_Hash value 334")
	assert.Nil(t, err)
}
