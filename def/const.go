package def

type AiRequest struct {
	Data        []interface{} `json:"data"`
	EventData   interface{}   `json:"event_data"`
	FnIndex     int           `json:"fn_index"`
	SessionHash string        `json:"session_hash"`
}

type AiResponse struct {
	Msg    string `json:"msg"`
	Output struct {
		Data         []interface{} `json:"data"`
		IsGenerating bool          `json:"is_generating"`
	} `json:"output"`
	Success bool `json:"success"`
}

type OpenAIChatRequest struct {
	Model    string              `json:"model"`
	Stream   bool                `json:"stream"`
	Messages []OpenAIChatMessage `json:"messages"`
}

type OpenAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIChatResponse struct {
	Choices []OpenAIChatChoice `json:"choices"`
}

type OpenAIChatChoice struct {
	Delta OpenAIChatDelta `json:"delta"`
}

type OpenAIChatDelta struct {
	Content string `json:"content"`
}
