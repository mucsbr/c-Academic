package s2s

import (
	"nixiang-gpt/def"
)

// 定义消息结构体
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 定义请求结构体
type ReqMessages struct {
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Model    string    `json:"model"`
}

//func ExtractConversations(messages []Message) [][]string {
//	var conversations [][]string
//	i := 0
//	for i < len(messages) {
//		if messages[i].Role == "user" {
//			userMessage := messages[i].Content
//			if i+1 < len(messages) && messages[i+1].Role == "assistant" {
//				assistantMessage := messages[i+1].Content
//				msgTemplate := "<div class=\"markdown-body\"><p>%s</p></div>"
//				tempUserMessage := fmt.Sprintf(msgTemplate, userMessage)
//				tempAssistantMessage := fmt.Sprintf(msgTemplate, assistantMessage)
//				conversations = append(conversations, []string{tempUserMessage, tempAssistantMessage})
//				i += 2
//			} else {
//				msgTemplate := "<div class=\"markdown-body\"><p>%s</p></div>"
//				tempUserMessage := fmt.Sprintf(msgTemplate, userMessage)
//				conversations = append(conversations, []string{tempUserMessage})
//				i++
//			}
//		} else {
//			i++ // 跳过system和未匹配的assistant消息
//		}
//	}
//	return conversations
//}

func ExtractConversations(messages []def.OpenAIChatMessage) [][]string {
	var conversations [][]string
	i := 0
	for i < len(messages) {
		if messages[i].Role == "user" {
			userMessage := messages[i].Content
			if i+1 < len(messages) && messages[i+1].Role == "assistant" {
				assistantMessage := messages[i+1].Content
				conversations = append(conversations, []string{userMessage, assistantMessage})
				i += 2
			} else {
				conversations = append(conversations, []string{userMessage})
				i++
			}
		} else {
			i++ // 跳过system和未匹配的assistant消息
		}
	}
	return conversations
}
