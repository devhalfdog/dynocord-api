package handlers

import "strings"

func removeBrackets(str string) string {
	rep := strings.NewReplacer(
		"[", "",
		"]", "",
	)

	return rep.Replace(str)
}
