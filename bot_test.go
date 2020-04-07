package probitpot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandF(t *testing.T) {
	min := 0.1
	max := 0.5
	r := randF(min, max)

	assert.GreaterOrEqual(t, r, min)
	assert.LessOrEqual(t, r, max)
}

func TestRandI(t *testing.T) {
	min := 1
	max := 5

	r := randI(min, max)

	assert.GreaterOrEqual(t, r, min)
	assert.LessOrEqual(t, r, max)
}

func TestRound(t *testing.T) {
	tests := []struct {
		number   float64
		expected float64
		p        int
	}{
		{
			4.3987376685,
			4.4,
			1,
		},
		{
			7.372635494,
			7.37,
			2,
		},
		{
			9.337636373,
			9.338,
			3,
		},
	}

	for _, test := range tests {
		r := round(test.number, test.p)

		assert.Equal(t, r, test.expected)
	}
}
