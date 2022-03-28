package module

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsLectureID(t *testing.T) {
	type expectOutputs struct {
		success bool
		id      int
	}
	patterns := []struct {
		expect expectOutputs // expectation
		given  string        // given input
	}{
		{expectOutputs{false, 0}, "12345"},
		{expectOutputs{false, 0}, "123456"},
		{expectOutputs{false, 0}, "#123456"},
		{expectOutputs{false, 0}, "#00000"},
		{expectOutputs{false, 0}, "#01234"},
		{expectOutputs{false, 0}, "1#12345"},
		{expectOutputs{false, 0}, "##12345"},
		{expectOutputs{false, 0}, "#1"},
		{expectOutputs{false, 0}, "#1234"},
		{expectOutputs{false, 0}, "＃12345"},
		{expectOutputs{false, 0}, "# 12345"},
		{expectOutputs{false, 0}, "#１１１１１"},
		{expectOutputs{true, 12345}, "#12345"},
		{expectOutputs{true, 99999}, "#99999"},
	}

	for _, p := range patterns {
		success, id := IsLectureID(p.given)
		assert.Equal(t, p.expect.success, success)
		assert.Equal(t, p.expect.id, id)
	}
}
