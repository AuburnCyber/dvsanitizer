package safety

import "strconv"

// If safe, convert the given string to a UInt32.
func StrToUInt32(input string) (retValue uint32, retConversionOK bool) {
	intVal, err := strconv.Atoi(input)
	if err != nil {
		return 0, false
	}
	return IntToUInt32(intVal)
}

// If safe, convert the given integer to a uint32.
func IntToUInt32(input int) (retValue uint32, retConversionOK bool) {
	switch {
	case input < 0:
		return 0, false
	case input > int(uint(^uint32(0))): // Probably too explicit but obviously safe.
		return 0, false
	}

	return uint32(input), true
}
