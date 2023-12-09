package model

import "regexp"



var TypeSuffix = []string{"t", "1s", "5s", "1m", "5m"}

var pattern = `\b\w+\.\w+\.\w+\.(t|1s|5s|1m|5m)\b`
var regex = regexp.MustCompile(pattern)

func IsSymbolValid(symbol string) bool {
	return regex.MatchString(symbol)
}