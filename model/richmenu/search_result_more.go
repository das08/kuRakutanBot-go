// This file was generated from JSON Schema using quicktype, do not modify it directly.

package model

import "encoding/json"

func UnmarshalSearchResultMore(data []byte) (SearchResultMore, error) {
	var r SearchResultMore
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SearchResultMore) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SearchResultMore struct {
	Type   string                 `json:"type"`
	Header SearchResultMoreHeader `json:"header"`
	Body   SearchResultMoreBody   `json:"body"`
	Styles Styles                 `json:"styles"`
}

type SearchResultMoreBody struct {
	Type     string                        `json:"type"`
	Layout   string                        `json:"layout"`
	Contents []SearchResultMoreBodyContent `json:"contents"`
}

type SearchResultMoreBodyContent struct {
	Type     string                         `json:"type"`
	Layout   *string                        `json:"layout,omitempty"`
	Contents []SearchResultMoreBodyContents `json:"contents,omitempty"`
	Margin   *string                        `json:"margin,omitempty"`
	Spacing  *string                        `json:"spacing,omitempty"`
}

type SearchResultMoreBodyContents struct {
	Type     string                 `json:"type"`
	Text     *string                `json:"text,omitempty"`
	Size     *string                `json:"size,omitempty"`
	Color    *string                `json:"color,omitempty"`
	Flex     *int64                 `json:"flex,omitempty"`
	Layout   *string                `json:"layout,omitempty"`
	Contents []SearchResultContents `json:"contents,omitempty"`
	Margin   *string                `json:"margin,omitempty"`
}

type SearchResultMoreHeader struct {
	Type          string                          `json:"type"`
	Layout        string                          `json:"layout"`
	Contents      []SearchResultMoreHeaderContent `json:"contents"`
	PaddingAll    string                          `json:"paddingAll"`
	Spacing       string                          `json:"spacing"`
	PaddingTop    string                          `json:"paddingTop"`
	PaddingBottom string                          `json:"paddingBottom"`
}

type SearchResultMoreHeaderContent struct {
	Type   string  `json:"type"`
	Text   *string `json:"text,omitempty"`
	Weight *string `json:"weight,omitempty"`
	Color  *string `json:"color,omitempty"`
	Size   *string `json:"size,omitempty"`
	Wrap   *bool   `json:"wrap,omitempty"`
}
