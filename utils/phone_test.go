package utils

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetPhoneByCountry(t *testing.T) {
	phoneList := []string{
		"700000000", "748000000", "110000000", "730000000", "762000000",
		"100000000", "101000000", "102000000",
		"112000000", "113000000", "114000000", "115000000",
		"779000000", "764000000", "747000000",
	}

	for _, number := range phoneList {
		t.Run(fmt.Sprint(number), func(t *testing.T) {
			phone, err := GetPhoneByCountry("KE", number)
			require.NoError(t, err)
			require.Equal(t, phone, "254"+number)
		})
	}
}
