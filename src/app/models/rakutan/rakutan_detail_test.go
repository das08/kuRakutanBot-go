package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLatestDetail(t *testing.T) {
	type expectOutputs struct {
		accept int
		total  int
	}
	patterns := []struct {
		expect expectOutputs  // expectation
		given  RakutanDetails // given input
	}{
		{expectOutputs{0, 0}, RakutanDetails{{2021, 0, 0}, {2020, 0, 0}, {2019, 0, 0}}},
		{expectOutputs{0, 255}, RakutanDetails{{2021, 0, 255}, {2020, 0, 0}, {2019, 0, 0}}},
		{expectOutputs{0, 255}, RakutanDetails{{2021, 0, 0}, {2020, 0, 255}, {2019, 0, 0}}},
		{expectOutputs{0, 255}, RakutanDetails{{2021, 0, 0}, {2020, 0, 0}, {2019, 0, 255}}},
		{expectOutputs{61, 120}, RakutanDetails{{2021, 61, 120}, {2020, 100, 100}, {2019, 70, 80}}},
		{expectOutputs{73, 120}, RakutanDetails{{2021, 73, 120}, {2020, 100, 100}, {2019, 70, 80}}},
		{expectOutputs{73, 100}, RakutanDetails{{2021, 73, 100}, {2020, 100, 100}, {2019, 70, 80}}},
		{expectOutputs{76, 100}, RakutanDetails{{2021, 76, 100}, {2020, 100, 100}, {2019, 70, 80}}},
		{expectOutputs{80, 100}, RakutanDetails{{2021, 80, 100}, {2020, 100, 100}, {2019, 70, 80}}},
		{expectOutputs{88, 100}, RakutanDetails{{2021, 88, 100}, {2020, 100, 100}, {2019, 70, 80}}},
		{expectOutputs{92, 100}, RakutanDetails{{2021, 92, 100}, {2020, 100, 100}, {2019, 70, 80}}},
		{expectOutputs{100, 100}, RakutanDetails{{2021, 0, 0}, {2020, 100, 100}, {2019, 70, 80}}},
	}

	for _, p := range patterns {
		accept, total := p.given.GetLatestDetail()
		assert.Equal(t, p.expect.accept, accept)
		assert.Equal(t, p.expect.total, total)
	}
}
