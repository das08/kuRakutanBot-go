// This file was generated from JSON Schema using quicktype, do not modify it directly.

package model

import "encoding/json"

func UnmarshalHelp(data []byte) (Help, error) {
	var r Help
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Help) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Help struct {
	Type   string     `json:"type"`
	Size   string     `json:"size"`
	Body   HelpBody   `json:"body"`
	Footer HelpFooter `json:"footer"`
	Styles Styles     `json:"styles"`
}

type HelpBody struct {
	Type     string            `json:"type"`
	Layout   string            `json:"layout"`
	Contents []HelpBodyContent `json:"contents"`
}

type HelpBodyContent struct {
	Type     string             `json:"type"`
	Text     *string            `json:"text,omitempty"`
	Weight   *string            `json:"weight,omitempty"`
	Size     *string            `json:"size,omitempty"`
	Align    *string            `json:"align,omitempty"`
	Margin   *string            `json:"margin,omitempty"`
	Layout   *string            `json:"layout,omitempty"`
	Spacing  *Spacing           `json:"spacing,omitempty"`
	Contents []HelpBodyContents `json:"contents,omitempty"`
}

type HelpBodyContents struct {
	Type     string        `json:"type"`
	Layout   string        `json:"layout"`
	Margin   string        `json:"margin"`
	Contents []HelpContent `json:"contents"`
}

type HelpContent struct {
	Type        Type        `json:"type"`
	Text        string      `json:"text"`
	Size        Spacing     `json:"size"`
	Align       Align       `json:"align"`
	Weight      *string     `json:"weight,omitempty"`
	Flex        *int64      `json:"flex,omitempty"`
	Action      *TextAction `json:"action,omitempty"`
	Color       *string     `json:"color,omitempty"`
	OffsetStart *string     `json:"offsetStart,omitempty"`
	Gravity     *string     `json:"gravity,omitempty"`
	Wrap        *bool       `json:"wrap,omitempty"`
}

type HelpFooter struct {
	Type     string              `json:"type"`
	Layout   string              `json:"layout"`
	Contents []HelpFooterContent `json:"contents"`
}

type HelpFooterContent struct {
	Type  Type   `json:"type"`
	Text  string `json:"text"`
	Wrap  bool   `json:"wrap"`
	Size  string `json:"size"`
	Color string `json:"color"`
}

// type Styles struct {
// 	Body   StylesBody   `json:"body"`
// 	Footer StylesFooter `json:"footer"`
// }

type StylesBody struct {
	BackgroundColor string `json:"backgroundColor"`
}

type StylesFooter struct {
	Separator       bool   `json:"separator"`
	BackgroundColor string `json:"backgroundColor"`
}

type Align string

const (
	Start Align = "start"
)

type Type string

const (
	Text Type = "text"
)
