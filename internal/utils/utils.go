package utils

import (
	"log"
	"strings"
)

func StringContains(s []string, str string) bool {
	for _, v := range s {
		if strings.Contains(v, str) {
			return true
		}
	}
	log.Println("String contains FALSE", s, str)
	return false
}
