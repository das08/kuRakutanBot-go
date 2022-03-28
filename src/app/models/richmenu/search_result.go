// This file was generated from JSON Schema using quicktype, do not modify it directly.
package richmenu

import "encoding/json"

func UnmarshalSearchResult(data []byte) (SearchResult, error) {
	var r SearchResult
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SearchResult) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SearchResult struct {
	Type   string             `json:"type"`
	Header SearchResultHeader `json:"header"`
	Body   SearchResultBody   `json:"body"`
	Styles Styles             `json:"styles"`
}

type SearchResultBody struct {
	Type     string        `json:"type"`
	Layout   string        `json:"layout"`
	Contents []BodyContent `json:"contents"`
}

type BodyContent struct {
	Type     string          `json:"type"`
	Layout   *string         `json:"layout,omitempty"`
	Contents []PurpleContent `json:"contents,omitempty"`
	Margin   *string         `json:"margin,omitempty"`
	Spacing  *string         `json:"spacing,omitempty"`
}

type PurpleContent struct {
	Type     string          `json:"type"`
	Text     *string         `json:"text,omitempty"`
	Size     *string         `json:"size,omitempty"`
	Color    *string         `json:"color,omitempty"`
	Flex     *int64          `json:"flex,omitempty"`
	Layout   *string         `json:"layout,omitempty"`
	Contents []FluffyContent `json:"contents,omitempty"`
	Margin   *string         `json:"margin,omitempty"`
}

type FluffyContent struct {
	Type         string  `json:"type"`
	Text         string  `json:"text"`
	Flex         int64   `json:"flex"`
	Size         string  `json:"size"`
	Weight       *string `json:"weight,omitempty"`
	Color        string  `json:"color"`
	Wrap         *bool   `json:"wrap,omitempty"`
	Align        *string `json:"align,omitempty"`
	Decoration   *string `json:"decoration,omitempty"`
	Margin       *string `json:"margin,omitempty"`
	Action       *Action `json:"action,omitempty"`
	OffsetBottom *string `json:"offsetBottom,omitempty"`
}

type Action struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	Text  string `json:"text"`
}

type SearchResultHeader struct {
	Type          string          `json:"type"`
	Layout        string          `json:"layout"`
	Contents      []HeaderContent `json:"contents"`
	PaddingAll    string          `json:"paddingAll"`
	Spacing       string          `json:"spacing"`
	PaddingTop    string          `json:"paddingTop"`
	PaddingBottom string          `json:"paddingBottom"`
}

type HeaderContent struct {
	Type   string  `json:"type"`
	Text   *string `json:"text,omitempty"`
	Weight *string `json:"weight,omitempty"`
	Color  *string `json:"color,omitempty"`
	Size   *string `json:"size,omitempty"`
	Wrap   *bool   `json:"wrap,omitempty"`
	Margin *string `json:"margin,omitempty"`
}

func (pc PurpleContent) DeepCopy() PurpleContent {
	tmp := pc
	tmp.Contents = make([]FluffyContent, len(pc.Contents))
	copy(tmp.Contents, pc.Contents)
	return tmp
}
