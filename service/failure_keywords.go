package service

import (
	"strings"
)

// CheckFailureKeywords 检测响应体是否包含失败关键字
// 返回 matched (true 表示匹配到关键字，即失败), matchedKeyword (匹配到的关键字)
func CheckFailureKeywords(responseBody string, keywords []string, caseSensitive bool) (bool, string) {
	if len(keywords) == 0 || responseBody == "" {
		return false, ""
	}

	bodyToCheck := responseBody
	if !caseSensitive {
		bodyToCheck = strings.ToLower(responseBody)
	}

	for _, keyword := range keywords {
		if keyword == "" {
			continue
		}

		keywordToCheck := keyword
		if !caseSensitive {
			keywordToCheck = strings.ToLower(keyword)
		}

		if strings.Contains(bodyToCheck, keywordToCheck) {
			return true, keyword
		}
	}

	return false, ""
}
