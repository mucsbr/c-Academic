package s2s

import (
	"encoding/json"
	"fmt"
	"nixiang-gpt/def"
	"testing"
)

func TestExtractConversations(t *testing.T) {
	reqMessagesJSON := `{
		"model": "gpt-4o",
		"stream": true,
		"messages": [
			{"role": "system", "content": "You are ChatGPT, a large language model trained by OpenAI. Knowledge cutoff: 2023-10 Current model: gpt-4o Current time: Wed Jul 17 2024 14:36:08 GMT+0800 (中国标准时间) Latex inline: \\(x^2\\) Latex block: $$e=mc^2$$"},
			{"role": "user", "content": "你好"},
			{"role": "assistant", "content": "你好！有什么我可以帮忙的吗？"},
			{"role": "user", "content": "鲁迅为什么打周树人"}
		]
	}`

	var reqMessages def.OpenAIChatRequest
	err := json.Unmarshal([]byte(reqMessagesJSON), &reqMessages)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// 提取并打印对话
	conversations := ExtractConversations(reqMessages.Messages)
	for _, conv := range conversations {
		fmt.Println(conv)
	}
}
