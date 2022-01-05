package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableDrivenCases(t *testing.T) {
	tests := map[string]struct {
		// other arguments...
		expectedError string
	}{
		// TODO: test cases
		"GivenCondition_WhenAction_ThenExpectResult": {},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange...

			// act...
			var err error
			main()
			// assert...
			if tt.expectedError != "" {
				require.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
