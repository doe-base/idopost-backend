package utils

import "strings"

func StringJsonFormatter(str []byte) string {
	s := string(str)
	t := strings.Replace(s, "\\n", "", -1)
	tt := strings.Replace(t, "\\", "", -1)
	ttt := strings.Replace(tt, " ", "", -1)
	tttt := strings.Replace(ttt, "\"{", "{", -1)
	ttttt := strings.Replace(tttt, "}\"", "}", -1)
	tttttt := strings.Replace(ttttt, "\"[", "[", -1)
	formatedString := strings.Replace(tttttt, "]\"", "]", -1)

	return formatedString
}
