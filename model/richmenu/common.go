package model

type Styles struct {
	Header HeaderClass `json:"header"`
	Body   HeaderClass `json:"body"`
	Footer Footer      `json:"footer"`
}

type HeaderClass struct {
	BackgroundColor string `json:"backgroundColor"`
}

type Footer struct {
	Separator       bool   `json:"separator"`
	BackgroundColor string `json:"backgroundColor"`
	SeparatorColor  string `json:"separatorColor"`
}

type TextAction struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	Text  string `json:"text"`
}

type URIAction struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	URI   string `json:"uri"`
}

type PostbackAction struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	Data  string `json:"data"`
}
