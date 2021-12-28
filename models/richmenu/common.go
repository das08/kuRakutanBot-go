package richmenu

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

type Spacing string

const (
	Sm  Spacing = "sm"
	Xs  Spacing = "xs"
	Xxs Spacing = "xxs"
)

type Type string

const (
	Text Type = "text"
)

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

type TextandPBAction struct {
	Type  string  `json:"type"`
	Label string  `json:"label"`
	Text  *string `json:"text,omitempty"`
	Data  *string `json:"data,omitempty"`
}
