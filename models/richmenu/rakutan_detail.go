// This file was generated from JSON Schema using quicktype, do not modify it directly.

package richmenu

import "encoding/json"

func UnmarshalRakutanDetail(data []byte) (RakutanDetail, error) {
	var r RakutanDetail
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *RakutanDetail) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type RakutanDetail struct {
	Type   string              `json:"type"`
	Header RakutanDetailHeader `json:"header"`
	Body   RakutanDetailBody   `json:"body"`
	Styles Styles              `json:"styles"`
}

type RakutanDetailBody struct {
	Type     BodyType                   `json:"type"`
	Layout   string                     `json:"layout"`
	Contents []RakutanDetailBodyContent `json:"contents"`
}

type RakutanDetailBodyContent struct {
	Type     BodyType               `json:"type"`
	Layout   *string                `json:"layout,omitempty"`
	Margin   string                 `json:"margin"`
	Spacing  *Spacing               `json:"spacing,omitempty"`
	Contents []RakutanDetailContent `json:"contents,omitempty"`
}

type RakutanDetailContent struct {
	Type     BodyType                `json:"type"`
	Layout   *string                 `json:"layout,omitempty"`
	Contents []RakutanDetailContents `json:"contents,omitempty"`
	Margin   *string                 `json:"margin,omitempty"`
	Text     *string                 `json:"text,omitempty"`
	Size     *Spacing                `json:"size,omitempty"`
	Color    *string                 `json:"color,omitempty"`
	Flex     *int64                  `json:"flex,omitempty"`
	Wrap     *bool                   `json:"wrap,omitempty"`
}

type RakutanDetailContents struct {
	Type        PurpleType `json:"type"`
	Text        *string    `json:"text,omitempty"`
	Size        *Spacing   `json:"size,omitempty"`
	Color       *string    `json:"color,omitempty"`
	Flex        *int64     `json:"flex,omitempty"`
	Align       *string    `json:"align,omitempty"`
	Style       *string    `json:"style,omitempty"`
	Weight      *string    `json:"weight,omitempty"`
	OffsetEnd   *string    `json:"offsetEnd,omitempty"`
	Wrap        *bool      `json:"wrap,omitempty"`
	Decoration  *string    `json:"decoration,omitempty"`
	Action      *URIAction `json:"action,omitempty"`
	URL         *string    `json:"url,omitempty"`
	AspectRatio *string    `json:"aspectRatio,omitempty"`
}

type RakutanDetailHeader struct {
	Type          BodyType                     `json:"type"`
	Layout        string                       `json:"layout"`
	Contents      []RakutanDetailHeaderContent `json:"contents"`
	PaddingAll    string                       `json:"paddingAll"`
	Spacing       string                       `json:"spacing"`
	PaddingTop    string                       `json:"paddingTop"`
	PaddingBottom string                       `json:"paddingBottom"`
}

type RakutanDetailHeaderContent struct {
	Type     BodyType                      `json:"type"`
	Layout   *string                       `json:"layout,omitempty"`
	Contents []RakutanDetailHeaderContents `json:"contents,omitempty"`
	Text     *string                       `json:"text,omitempty"`
	Weight   *string                       `json:"weight,omitempty"`
	Size     *string                       `json:"size,omitempty"`
	Margin   *string                       `json:"margin,omitempty"`
	Color    *string                       `json:"color,omitempty"`
	Wrap     *bool                         `json:"wrap,omitempty"`
	Spacing  *Spacing                      `json:"spacing,omitempty"`
}

type RakutanDetailHeaderContents struct {
	Type        PurpleType      `json:"type"`
	URL         *string         `json:"url,omitempty"`
	AspectRatio *string         `json:"aspectRatio,omitempty"`
	Flex        *int64          `json:"flex,omitempty"`
	OffsetStart *string         `json:"offsetStart,omitempty"`
	Action      *PostbackAction `json:"action,omitempty"`
	Text        *string         `json:"text,omitempty"`
	Weight      *string         `json:"weight,omitempty"`
	Color       *string         `json:"color,omitempty"`
	Size        *Spacing        `json:"size,omitempty"`
	Align       *string         `json:"align,omitempty"`
}

type PurpleType string

const (
	Filler     PurpleType = "filler"
	Image      PurpleType = "image"
	PurpleText PurpleType = "text"
)

type BodyType string

const (
	Box        BodyType = "box"
	FluffyText BodyType = "text"
	Separator  BodyType = "separator"
	Spacer     BodyType = "spacer"
)
