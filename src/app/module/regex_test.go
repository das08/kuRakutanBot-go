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

func TestParsePBParam(t *testing.T) {
	type expectOutputs struct {
		success bool
		params  PostbackParam
	}
	patterns := []struct {
		expect expectOutputs // expectation
		given  string        // given input
	}{
		{expectOutputs{false, PostbackParam{}}, "type=fav&id=12345"},
		{expectOutputs{false, PostbackParam{}}, "type=fav&lecname=sample&id=12345"},
		{expectOutputs{false, PostbackParam{}}, "id=12345&type=fav&lecname=sample"},
		{expectOutputs{false, PostbackParam{}}, "types=fav&id=12345&lecname=sample"},
		{expectOutputs{false, PostbackParam{}}, "type=fav&ids=12345&lecname=sample"},
		{expectOutputs{false, PostbackParam{}}, "type=fav&id=12345&lecnames=sample"},
		{expectOutputs{false, PostbackParam{}}, "TYPE=fav&id=12345&lecname=sample"},
		{expectOutputs{false, PostbackParam{}}, "type=fav&ID=12345&lecname=sample"},
		{expectOutputs{false, PostbackParam{}}, "type=fav&id=12345&LECNAME=sample"},
		{expectOutputs{false, PostbackParam{}}, ""},
		{expectOutputs{false, PostbackParam{}}, "=&=&="},
		{expectOutputs{true, PostbackParam{"fav", 12345, "sample"}}, "type=fav&id=12345&lecname=sample"},
		{expectOutputs{true, PostbackParam{"fav", 99999, "日本語"}}, "type=fav&id=99999&lecname=日本語"},
	}

	for _, p := range patterns {
		success, params := ParsePBParam(p.given)
		assert.Equal(t, p.expect.success, success)
		assert.Equal(t, p.expect.params.Type, params.Type)
		assert.Equal(t, p.expect.params.ID, params.ID)
		assert.Equal(t, p.expect.params.LectureName, params.LectureName)
	}
}
