package service

import (
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
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

// CheckFailureKeywordsInStreamContent 检测流式响应的 content 字段是否包含失败关键字
// 会尝试解析 JSON 并提取 choices[].delta.content 字段进行检测
// 返回 matched (true 表示匹配到关键字，即失败), matchedKeyword (匹配到的关键字)
func CheckFailureKeywordsInStreamContent(data string, keywords []string, caseSensitive bool) (bool, string) {
	if len(keywords) == 0 || data == "" {
		return false, ""
	}

	// 尝试解析为 ChatCompletionsStreamResponse
	var streamResponse dto.ChatCompletionsStreamResponse
	if err := common.UnmarshalJsonStr(data, &streamResponse); err != nil {
		// 解析失败，不进行检测（可能不是有效的流式响应）
		return false, ""
	}

	// 提取所有 content 和 reasoning_content 进行检测
	var contentBuilder strings.Builder
	for _, choice := range streamResponse.Choices {
		contentBuilder.WriteString(choice.Delta.GetContentString())
		contentBuilder.WriteString(choice.Delta.GetReasoningContent())
	}

	content := contentBuilder.String()
	if content == "" {
		// 没有 content，不进行检测
		return false, ""
	}

	return CheckFailureKeywords(content, keywords, caseSensitive)
}

// CheckFailureKeywordsInClaudeContent 检测 Claude 流式响应的 content 字段是否包含失败关键字
// 会尝试解析 JSON 并提取 delta.text 或 content_block.text 字段进行检测
// 返回 matched (true 表示匹配到关键字，即失败), matchedKeyword (匹配到的关键字)
func CheckFailureKeywordsInClaudeContent(data string, keywords []string, caseSensitive bool) (bool, string) {
	if len(keywords) == 0 || data == "" {
		return false, ""
	}

	// 尝试解析为 ClaudeResponse
	var claudeResponse dto.ClaudeResponse
	if err := common.UnmarshalJsonStr(data, &claudeResponse); err != nil {
		// 解析失败，不进行检测
		return false, ""
	}

	// 提取 delta.text 或其他文本内容进行检测
	var contentBuilder strings.Builder

	// 检查 delta 中的文本
	if claudeResponse.Delta != nil {
		if claudeResponse.Delta.Text != nil {
			contentBuilder.WriteString(*claudeResponse.Delta.Text)
		}
		if claudeResponse.Delta.Thinking != nil {
			contentBuilder.WriteString(*claudeResponse.Delta.Thinking)
		}
	}

	// 检查 content_block 中的文本
	if claudeResponse.ContentBlock != nil {
		if claudeResponse.ContentBlock.Text != nil {
			contentBuilder.WriteString(*claudeResponse.ContentBlock.Text)
		}
	}

	// 检查 content 数组中的文本
	for _, content := range claudeResponse.Content {
		if content.Text != nil {
			contentBuilder.WriteString(*content.Text)
		}
	}

	content := contentBuilder.String()
	if content == "" {
		// 没有 content，不进行检测
		return false, ""
	}

	return CheckFailureKeywords(content, keywords, caseSensitive)
}
// 会尝试解析 JSON 并提取 candidates[].content.parts[].text 字段进行检测
// 返回 matched (true 表示匹配到关键字，即失败), matchedKeyword (匹配到的关键字)
func CheckFailureKeywordsInGeminiContent(data string, keywords []string, caseSensitive bool) (bool, string) {
	if len(keywords) == 0 || data == "" {
		return false, ""
	}

	// 尝试解析为 GeminiChatResponse
	var geminiResponse dto.GeminiChatResponse
	if err := common.UnmarshalJsonStr(data, &geminiResponse); err != nil {
		// 解析失败，不进行检测
		return false, ""
	}

	// 提取所有 text 进行检测
	var contentBuilder strings.Builder
	for _, candidate := range geminiResponse.Candidates {
		for _, part := range candidate.Content.Parts {
			if part.Text != "" {
				contentBuilder.WriteString(part.Text)
			}
		}
	}

	content := contentBuilder.String()
	if content == "" {
		// 没有 content，不进行检测
		return false, ""
	}

	return CheckFailureKeywords(content, keywords, caseSensitive)
}
