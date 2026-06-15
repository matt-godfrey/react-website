package quotes

type Quote struct {
	Text      string `json:"text"`
	Author    string `json:"author"`
	CharCount any    `json:"charCount"`
	Html      string `json:"html"`
}
