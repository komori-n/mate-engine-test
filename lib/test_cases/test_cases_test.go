package test_cases

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestSet_TimeLimit(t *testing.T) {
	ts, err := Decode(`---
hoge:
    time_limit: 10`)
	assert.Nil(t, err)

	assert.Equal(t, ts["hoge"].TimeLimit, 10)
}

func TestTestSet_Opts(t *testing.T) {
	ts, err := Decode(`---
hoge:
    engine_opts:
        A: 334
        B: 264
`)
	assert.Nil(t, err)

	assert.Equal(t, ts["hoge"].Opts["A"], "334")
	assert.Equal(t, ts["hoge"].Opts["B"], "264")
}

func TestTestSet_NoOpts(t *testing.T) {
	ts, err := Decode(`---
hoge:
`)
	assert.Nil(t, err)

	assert.Equal(t, len(ts["hoge"].Opts), 0)
}

func TestTestSet_TestsSfen(t *testing.T) {
	ts, err := Decode(`---
hoge:
    tests:
        - sfen: 334
`)
	assert.Nil(t, err)

	_ = ts

	assert.Equal(t, len(ts["hoge"].Tests), 1)
	assert.Equal(t, ts["hoge"].Tests[0].Sfen, "334")
	assert.False(t, ts["hoge"].Tests[0].NoMate)
}

func TestTestSet_TestsNoMate(t *testing.T) {
	ts, err := Decode(`---
hoge:
    tests:
        - sfen: 264
          nomate: true
`)
	assert.Nil(t, err)

	_ = ts

	assert.Equal(t, len(ts["hoge"].Tests), 1)
	assert.Equal(t, ts["hoge"].Tests[0].Sfen, "264")
	assert.True(t, ts["hoge"].Tests[0].NoMate)
}
