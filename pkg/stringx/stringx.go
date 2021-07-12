package stringx

import "strings"

// FilterEmpty filter empty items
func FilterEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// SplitNoEmpty split string with seperator and filter empty items
func SplitNoEmpty(s string, sep string) []string {
	return FilterEmpty(strings.Split(s, sep))
}
