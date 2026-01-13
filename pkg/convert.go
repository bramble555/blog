package pkg

import (
	"strconv"
)

// StringSliceToInt64Slice converts a slice of strings to a slice of int64.
func StringSliceToInt64Slice(s []string) ([]int64, error) {
	res := make([]int64, len(s))
	for i, v := range s {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		res[i] = n
	}
	return res, nil
}
