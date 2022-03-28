package module

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsLectureNumber(t *testing.T) {
	type expectOutputs struct {
		success bool
		id      int
	}
	patterns := []struct {
		expect expectOutputs // expectation
		given  string        // given input
	}{
		{expectOutputs{false, 0}, "123456"},
		{expectOutputs{false, 0}, "1234567"},
		{expectOutputs{false, 0}, "#1234567"},
		{expectOutputs{false, 0}, "#000000"},
		{expectOutputs{false, 0}, "#012345"},
		{expectOutputs{false, 0}, "1#123456"},
		{expectOutputs{false, 0}, "##123456"},
		{expectOutputs{false, 0}, "#1"},
		{expectOutputs{false, 0}, "#12345"},
		{expectOutputs{false, 0}, "＃123456"},
		{expectOutputs{false, 0}, "# 123456"},
		{expectOutputs{false, 0}, "#１１１１１１"},
		{expectOutputs{true, 123456}, "#123456"},
		{expectOutputs{true, 999999}, "#999999"},
	}

	for _, p := range patterns {
		success, id := IsLectureNumber(p.given)
		assert.Equal(t, p.expect.success, success)
		assert.Equal(t, p.expect.id, id)
	}
}
