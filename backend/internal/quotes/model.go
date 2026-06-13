package quotes

type Quote struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Author    string `json:"author"`
	CharCount int    `json:"charCount"`
	Html      string `json:"html"`
}
