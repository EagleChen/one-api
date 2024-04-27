package cloudflare

type chatResponse struct {
	Result   chatCompleteResponse `json:"result"`
	Success  bool                 `json:"success"`
	Errors   []string             `json:"errors"`
	Messages []string             `json:"messages"`
}

type chatCompleteResponse struct {
	Response string `json:"response"`
}

type streamResponse = chatCompleteResponse
