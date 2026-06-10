package utils

import "strconv"

func StringToInt64(
	value string,
) int64 {

	result, _ := strconv.ParseInt(
		value,
		10,
		64,
	)

	return result
}
