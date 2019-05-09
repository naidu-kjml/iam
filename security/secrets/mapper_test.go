package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTokenName(t *testing.T) {
	mapper := Mapper{
		Configuration{
			[]ServiceConfiguration{
				{
					"balkan",
					[]string{"production", "simondev"},
					"BALKAN",
				},
				{
					"TRANSACTIONAL_MESSAGING_-_NEST_APP",
					[]string{"production"},
					"NEST",
				},
			},
		},
	}

	tests := [][]string{
		{"BALKAN_PRODUCTION", "Balkan", "Production"},
		{"BALKAN", "balkan", ""},
		{"NEST_PRODUCTION", "TRANSACTIONAL_MESSAGING_-_NEST_APP", "production"},
		{"", "Balkan", "notsimondev"},
		{"", "unknown", "whatever"},
	}

	for _, test := range tests {
		result, err := mapper.GetTokenName(test[1], test[2])
		if test[0] == "" {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test[0], result)
		}
	}

}
