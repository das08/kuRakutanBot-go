// This file was generated from JSON Schema using quicktype, do not modify it directly.

package richmenu

import "encoding/json"

func UnmarshalJudgeDetail(data []byte) (JudgeDetail, error) {
	var r JudgeDetail
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *JudgeDetail) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type JudgeDetail struct {
	Type   string            `json:"type"`
	Size   string            `json:"size"`
	Body   JudgeDetailBody   `json:"body"`
	Footer JudgeDetailFooter `json:"footer"`
	Styles JudgeDetailStyles `json:"styles"`
}

type JudgeDetailBody struct {
	Type     string                   `json:"type"`
	Layout   string                   `json:"layout"`
	Contents []JudgeDetailBodyContent `json:"contents"`
}

type JudgeDetailBodyContent struct {
	Type     string               `json:"type"`
	Text     *string              `json:"text,omitempty"`
	Weight   *Weight              `json:"weight,omitempty"`
	Size     *string              `json:"size,omitempty"`
	Align    *Align               `json:"align,omitempty"`
	Margin   *string              `json:"margin,omitempty"`
	Layout   *string              `json:"layout,omitempty"`
	Spacing  *Spacing             `json:"spacing,omitempty"`
	Contents []JudgeDetailContent `json:"contents,omitempty"`
}

type JudgeDetailContent struct {
	Type     string                `json:"type"`
	Layout   string                `json:"layout"`
	Margin   string                `json:"margin"`
	Contents []JudgeDetailContent2 `json:"contents"`
	Spacing  *string               `json:"spacing,omitempty"`
}

type JudgeDetailContent2 struct {
	Type   Type    `json:"type"`
	Text   string  `json:"text"`
	Size   Spacing `json:"size"`
	Color  string  `json:"color"`
	Style  *string `json:"style,omitempty"`
	Align  Align   `json:"align"`
	Weight *Weight `json:"weight,omitempty"`
}

type JudgeDetailFooter struct {
	Type     string          `json:"type"`
	Layout   string          `json:"layout"`
	Contents []FooterContent `json:"contents"`
}

type FooterContent struct {
	Type  Type   `json:"type"`
	Text  string `json:"text"`
	Color string `json:"color"`
	Wrap  bool   `json:"wrap"`
	Size  string `json:"size"`
}

type JudgeDetailStyles struct {
	Body   StylesBody   `json:"body"`
	Footer StylesFooter `json:"footer"`
}

const (
	Center Align = "center"
)

type Weight string

const (
	Bold Weight = "bold"
)
