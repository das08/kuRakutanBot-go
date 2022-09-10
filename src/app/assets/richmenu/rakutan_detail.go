package richmenu

import "github.com/line/line-bot-sdk-go/v7/linebot"

func LoadRakutanDetail() *linebot.BubbleContainer {
	container := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Contents: []linebot.FlexComponent{
						&linebot.ImageComponent{
							Type:        linebot.FlexComponentTypeImage,
							URL:         "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gray_star_28.png",
							AspectRatio: linebot.FlexImageAspectRatioType2to1,
							Flex:        toIntPtr(1),
							OffsetStart: "-5px",
							Action: &linebot.PostbackAction{
								Label: "action",
								Data:  "type=fav&id=12345&lecname=sample",
							},
						},
						&linebot.TextComponent{
							Type:        linebot.FlexComponentTypeText,
							Text:        "お気に入りに追加",
							Weight:      linebot.FlexTextWeightTypeBold,
							Color:       "#01E550",
							Size:        linebot.FlexTextSizeTypeSm,
							Flex:        toIntPtr(7),
							OffsetStart: "-5px",
						},
					},
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "[Lecture Name]",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeXxl,
					Margin: linebot.FlexComponentMarginTypeMd,
					Color:  "#E6E7E2",
					Wrap:   true,
				},
				&linebot.SeparatorComponent{
					Type: linebot.FlexComponentTypeSeparator,
				},
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeBaseline,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   "開講部局",
							Size:   linebot.FlexTextSizeTypeXs,
							Color:  "#81817F",
							Align:  linebot.FlexComponentAlignTypeStart,
							Flex:   toIntPtr(1),
							Weight: linebot.FlexTextWeightTypeBold,
						},
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   "[Faculty Name]",
							Size:   linebot.FlexTextSizeTypeXs,
							Color:  "#E6E7E2",
							Align:  linebot.FlexComponentAlignTypeStart,
							Flex:   toIntPtr(3),
							Weight: linebot.FlexTextWeightTypeBold,
						},
					},
					Spacing: linebot.FlexComponentSpacingTypeSm,
				},
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeBaseline,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   "群",
							Size:   linebot.FlexTextSizeTypeXs,
							Color:  "#81817F",
							Weight: linebot.FlexTextWeightTypeBold,
						},
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   "[Group Name]",
							Size:   linebot.FlexTextSizeTypeXs,
							Color:  "#E6E7E2",
							Weight: linebot.FlexTextWeightTypeBold,
						},
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   "単位数",
							Size:   linebot.FlexTextSizeTypeXs,
							Color:  "#81817F",
							Weight: linebot.FlexTextWeightTypeBold,
						},
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   "[Credit(s)]",
							Size:   linebot.FlexTextSizeTypeXs,
							Color:  "#E6E7E2",
							Weight: linebot.FlexTextWeightTypeBold,
						},
					},
				},
				&linebot.SpacerComponent{
					Type: linebot.FlexComponentTypeSpacer,
					Size: linebot.FlexSpacerSizeTypeLg,
				},
			},
			PaddingAll:    "20px",
			Spacing:       linebot.FlexComponentSpacingTypeMd,
			PaddingTop:    "22px",
			PaddingBottom: "0px",
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeVertical,
					Margin:  linebot.FlexComponentMarginTypeXxl,
					Spacing: linebot.FlexComponentSpacingTypeSm,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeHorizontal,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "単位取得率",
									Size:  linebot.FlexTextSizeTypeXxs,
									Color: "#aaaaaa",
									Flex:  toIntPtr(0),
								},
							},
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeHorizontal,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "2019年度",
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#555555",
									Flex:  toIntPtr(0),
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "0%",
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#111111",
									Align: linebot.FlexComponentAlignTypeEnd,
								},
							},
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeHorizontal,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "2018年度",
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#555555",
									Flex:  toIntPtr(0),
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "0%",
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#111111",
									Align: linebot.FlexComponentAlignTypeEnd,
								},
							},
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeHorizontal,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "2017年度",
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#555555",
									Flex:  toIntPtr(0),
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "0%",
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#111111",
									Align: linebot.FlexComponentAlignTypeEnd,
								},
							},
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeHorizontal,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "2016年度",
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#555555",
									Flex:  toIntPtr(0),
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "0%",
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#111111",
									Align: linebot.FlexComponentAlignTypeEnd,
								},
							},
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeHorizontal,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "2015年度",
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#555555",
									Flex:  toIntPtr(0),
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "0%",
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#111111",
									Align: linebot.FlexComponentAlignTypeEnd,
								},
							},
						},
						&linebot.SeparatorComponent{
							Type:   linebot.FlexComponentTypeSeparator,
							Margin: linebot.FlexComponentMarginTypeXxl,
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeHorizontal,
							Margin: linebot.FlexComponentMarginTypeXxl,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "らくたん判定",
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#555555",
								},
								&linebot.TextComponent{
									Type:      linebot.FlexComponentTypeText,
									Text:      "SSS",
									Size:      linebot.FlexTextSizeTypeSm,
									Color:     "#c3c45b",
									Align:     linebot.FlexComponentAlignTypeEnd,
									Style:     linebot.FlexTextStyleTypeItalic,
									Weight:    linebot.FlexTextWeightTypeBold,
									OffsetEnd: "6px",
								},
							},
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeBaseline,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "過去問",
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#555555",
									Flex:  toIntPtr(4),
								},
								&linebot.TextComponent{
									Type:   linebot.FlexComponentTypeText,
									Text:   "△",
									Weight: linebot.FlexTextWeightTypeBold,
									Align:  linebot.FlexComponentAlignTypeEnd,
									Color:  "#ffb101",
									Flex:   toIntPtr(4),
								},
								&linebot.TextComponent{
									Type:       linebot.FlexComponentTypeText,
									Text:       "追加する",
									Size:       linebot.FlexTextSizeTypeSm,
									Align:      linebot.FlexComponentAlignTypeEnd,
									Flex:       toIntPtr(3),
									Color:      "#777777",
									Decoration: linebot.FlexTextDecorationTypeUnderline,
									Action: &linebot.URIAction{
										Label: "action",
										URI:   "https://sites.google.com/view/siketai/%E9%81%8E%E5%8E%BB%E5%95%8F",
									},
								},
							},
						},
					},
				},
				&linebot.SeparatorComponent{
					Type:   linebot.FlexComponentTypeSeparator,
					Margin: linebot.FlexComponentMarginTypeXxl,
				},
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Margin: linebot.FlexComponentMarginTypeMd,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  "※単位取得率は「各授業の平均学生在籍数」をもとに表示しています。",
							Size:  linebot.FlexTextSizeTypeXxs,
							Color: "#aaaaaa",
							Flex:  toIntPtr(0),
							Wrap:  true,
						},
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  "※らくたん判定は単位取得率をもとに判定しています。詳しくは「判定詳細」と送信してください。",
							Size:  linebot.FlexTextSizeTypeXxs,
							Color: "#aaaaaa",
							Flex:  toIntPtr(0),
							Wrap:  true,
						},
					},
				},
			},
		},
		Styles: &linebot.BubbleStyle{
			Header: &linebot.BlockStyle{
				BackgroundColor: "#3C3C3A",
			},
			Body: &linebot.BlockStyle{
				BackgroundColor: "#F4F3F9",
			},
			Footer: &linebot.BlockStyle{
				Separator:       true,
				BackgroundColor: "#E7E8E3",
				SeparatorColor:  "#D6D6D4",
			},
		},
	}
	return container
}

func toIntPtr(i int) *int {
	return &i
}
