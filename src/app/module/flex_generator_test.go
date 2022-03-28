package module

import (
	models "github.com/das08/kuRakutanBot-go/models/rakutan"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRakutanPercent(t *testing.T) {
	type inputs struct {
		arg1 int
		arg2 int
	}
	patterns := []struct {
		expect string // expectation
		given  inputs // given input
	}{
		{"---% (0/0)", inputs{0, 0}},
		{"0.0% (0/5)", inputs{0, 5}},
		{"60.0% (3/5)", inputs{3, 5}},
		{"37.5% (3/8)", inputs{3, 8}},
		{"34.9% (117/335)", inputs{117, 335}},
		{"100.0% (335/335)", inputs{335, 335}},
	}

	for _, p := range patterns {
		actual := getRakutanPercent(p.given.arg1, p.given.arg2)
		assert.Equal(t, p.expect, actual)
	}
}

func TestGetRakutanJudge(t *testing.T) {
	patterns := []struct {
		expect RakutanJudge          // expectation
		given  models.RakutanDetails // given input
	}{
		{judgeList[8], models.RakutanDetails{{2021, 0, 0}, {2020, 0, 0}, {2019, 0, 0}}},
		{judgeList[7], models.RakutanDetails{{2021, 0, 255}, {2020, 0, 0}, {2019, 0, 0}}},
		{judgeList[7], models.RakutanDetails{{2021, 0, 0}, {2020, 0, 255}, {2019, 0, 0}}},
		{judgeList[7], models.RakutanDetails{{2021, 0, 0}, {2020, 0, 0}, {2019, 0, 255}}},
		{judgeList[6], models.RakutanDetails{{2021, 61, 120}, {2020, 100, 100}, {2019, 70, 80}}},
		{judgeList[5], models.RakutanDetails{{2021, 73, 120}, {2020, 100, 100}, {2019, 70, 80}}},
		{judgeList[4], models.RakutanDetails{{2021, 73, 100}, {2020, 100, 100}, {2019, 70, 80}}},
		{judgeList[3], models.RakutanDetails{{2021, 76, 100}, {2020, 100, 100}, {2019, 70, 80}}},
		{judgeList[2], models.RakutanDetails{{2021, 80, 100}, {2020, 100, 100}, {2019, 70, 80}}},
		{judgeList[1], models.RakutanDetails{{2021, 88, 100}, {2020, 100, 100}, {2019, 70, 80}}},
		{judgeList[0], models.RakutanDetails{{2021, 92, 100}, {2020, 100, 100}, {2019, 70, 80}}},
		{judgeList[0], models.RakutanDetails{{2021, 0, 0}, {2020, 100, 100}, {2019, 70, 80}}},
	}

	for _, p := range patterns {
		actual := getRakutanJudge(p.given)
		assert.Equal(t, p.expect, actual)
	}
}

func TestToStr(t *testing.T) {
	patterns := []struct {
		expect string // expectation
		given  int    // given input
	}{
		{"0", 0},
		{"128", 128},
		{"1234567890", 1234567890},
	}

	for _, p := range patterns {
		actual := toStr(p.given)
		assert.Equal(t, p.expect, actual)
	}
}

func TestToPtr(t *testing.T) {
	sampleText := "Hello World"
	expect := &sampleText
	actual := toPtr(sampleText)
	assert.Equal(t, expect, actual)
}
