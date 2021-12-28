// This file was generated from JSON Schema using quicktype, do not modify it directly.
package model

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
	Type     string                     `json:"type"`
	Layout   string                     `json:"layout"`
	Contents []RakutanDetailBodyContent `json:"contents"`
}

type BodyContent struct {
	Type     string                `json:"type"`
	Layout   *string               `json:"layout,omitempty"`
	Contents []SearchResultContent `json:"contents,omitempty"`
	Margin   *string               `json:"margin,omitempty"`
	Spacing  *string               `json:"spacing,omitempty"`
}

type SearchResultContent struct {
	Type     string                 `json:"type"`
	Text     *string                `json:"text,omitempty"`
	Size     *string                `json:"size,omitempty"`
	Color    *string                `json:"color,omitempty"`
	Flex     *int64                 `json:"flex,omitempty"`
	Layout   *string                `json:"layout,omitempty"`
	Contents []SearchResultContents `json:"contents,omitempty"`
	Margin   *string                `json:"margin,omitempty"`
}

type SearchResultContents struct {
	Type         string      `json:"type"`
	Text         string      `json:"text"`
	Size         string      `json:"size"`
	Color        string      `json:"color"`
	Flex         int64       `json:"flex"`
	Wrap         *bool       `json:"wrap,omitempty"`
	Align        *string     `json:"align,omitempty"`
	Weight       *string     `json:"weight,omitempty"`
	Decoration   *string     `json:"decoration,omitempty"`
	Margin       *string     `json:"margin,omitempty"`
	Action       *TextAction `json:"action,omitempty"`
	OffsetBottom *string     `json:"offsetBottom,omitempty"`
}

type SearchResultHeader struct {
	Type          string                      `json:"type"`
	Layout        string                      `json:"layout"`
	Contents      []SearchResultHeaderContent `json:"contents"`
	PaddingAll    string                      `json:"paddingAll"`
	Spacing       string                      `json:"spacing"`
	PaddingTop    string                      `json:"paddingTop"`
	PaddingBottom string                      `json:"paddingBottom"`
}

type SearchResultHeaderContent struct {
	Type   string  `json:"type"`
	Text   *string `json:"text,omitempty"`
	Weight *string `json:"weight,omitempty"`
	Color  *string `json:"color,omitempty"`
	Size   *string `json:"size,omitempty"`
	Wrap   *bool   `json:"wrap,omitempty"`
	Margin *string `json:"margin,omitempty"`
}
