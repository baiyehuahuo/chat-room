package logic

import (
	"chatroom/global"
	"strings"
)

func FilterSensitive(content string) string {
	for _, word := range global.SensitiveWords {
		content = strings.ReplaceAll(content, word, strings.Repeat("*", len([]rune(word))))
	}
	return content
}
