// This file was generated from JSON Schema using quicktype, do not modify it directly.

package richmenu

import "encoding/json"

func UnmarshalFavorites(data []byte) (Favorites, error) {
	var r Favorites
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Favorites) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Favorites struct {
	Type   string        `json:"type"`
	Header Header        `json:"header"`
	Body   FavoritesBody `json:"body"`
	Styles Styles        `json:"styles"`
}

type FavoritesBody struct {
	Type     string                 `json:"type"`
	Layout   string                 `json:"layout"`
	Contents []FavoritesBodyContent `json:"contents"`
}

type FavoritesBodyContent struct {
	Type     string                  `json:"type"`
	Layout   *string                 `json:"layout,omitempty"`
	Margin   string                  `json:"margin"`
	Spacing  *string                 `json:"spacing,omitempty"`
	Contents []FavoritesBodyContents `json:"contents,omitempty"`
}

type FavoritesBodyContents struct {
	Type     string             `json:"type"`
	Layout   *string            `json:"layout,omitempty"`
	Contents []FavoritesContent `json:"contents,omitempty"`
	Margin   *string            `json:"margin,omitempty"`
	Text     *string            `json:"text,omitempty"`
	Size     *string            `json:"size,omitempty"`
	Color    *string            `json:"color,omitempty"`
}

type FavoritesContent struct {
	Type         string           `json:"type"`
	Text         string           `json:"text"`
	Color        string           `json:"color"`
	Flex         *int64           `json:"flex,omitempty"`
	Wrap         *bool            `json:"wrap,omitempty"`
	Size         string           `json:"size"`
	Align        *string          `json:"align,omitempty"`
	Weight       *string          `json:"weight,omitempty"`
	Decoration   *string          `json:"decoration,omitempty"`
	Margin       *string          `json:"margin,omitempty"`
	Action       *TextandPBAction `json:"action,omitempty"`
	OffsetBottom *string          `json:"offsetBottom,omitempty"`
}

type Header struct {
	Type     string                   `json:"type"`
	Layout   string                   `json:"layout"`
	Contents []FavoritesHeaderContent `json:"contents"`
	Spacing  string                   `json:"spacing"`
}

type FavoritesHeaderContent struct {
	Type   string  `json:"type"`
	Text   string  `json:"text"`
	Weight string  `json:"weight"`
	Color  string  `json:"color"`
	Wrap   *bool   `json:"wrap,omitempty"`
	Align  *string `json:"align,omitempty"`
}
