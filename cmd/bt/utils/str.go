package utils

func StrPadLeft(input string, padLength int, padString string) string {
	output := ""
	inputLen := len(input)
	if inputLen >= padLength {
		return input
	}
	ll := padLength - inputLen
	for i := 1; i <= ll; i = i + len(padString) {
		output += padString
	}
	return output + input
}

func StrPadRight(input string, padLength int, padString string) string {
	output := ""
	inputLen := len(input)
	if inputLen >= padLength {
		return input
	}
	ll := padLength - inputLen
	for i := 1; i <= ll; i = i + len(padString) {
		output += padString
	}
	return input + output
}
