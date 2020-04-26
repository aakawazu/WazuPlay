package db

import "regexp"

// EscapeSinglequotation escape single quotation
func EscapeSinglequotation(str string) string {
	rep := regexp.MustCompile(`'`)
	str = rep.ReplaceAllString(str, "''")
	return str
}
