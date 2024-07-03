package teleblog

import (
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strings"
	"unicode/utf16"

	"github.com/pocketbase/pocketbase/tools/types"
	"gopkg.in/telebot.v3"
)

type MarkupNyPosition struct {
	Offset   int
	Priority int
	IsOpen   bool
	Tag      []rune
}

func FormHistoryTextWithMarkup(markup types.JsonArray[HistoryMessageTextEntity]) string {
	text := ""

	for _, entity := range markup {
		switch entity.Type {
		case telebot.EntityItalic:
			text += "<i class='inline'>" + entity.Text + "</i>"
		case telebot.EntityBold:
			text += "<b class='inline'>" + entity.Text + "</b>"
		case telebot.EntityURL:
			link := entity.Text
			if strings.Contains(link, "://") == false {
				link = "http://" + link
			}
			text += "<a target='_blank' href='" + link + "' class='inline c-link'>" + entity.Text + "</a>"
		case telebot.EntityTextLink:
			text += "<a target='_blank' class='inline c-link' href='" + entity.Text + "'>" + entity.Text + "</a>"
		case telebot.EntityMention:
			text += "<a target='_blank' href='https://t.me/" + entity.Text + "' class='inline c-link'>" + entity.Text + "</a>"
		default:
			text += entity.Text
		}
	}

	return text
}

func AddMarkupToText(srcText string, markup types.JsonArray[telebot.MessageEntity]) (string, error) {
	text := utf16.Encode([]rune(srcText))

	var entities telebot.Entities

	b, err := markup.MarshalJSON()
	if err != nil {
		return "", fmt.Errorf("failed to marshal markup: %w", err)
	}

	err = json.Unmarshal(b, &entities)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal markup: %w", err)
	}

	var markUpByPosition []MarkupNyPosition

	for i, entity := range entities {
		switch entity.Type {
		case telebot.EntityItalic:
			markUpByPosition = append(markUpByPosition, MarkupNyPosition{Offset: entity.Offset, Tag: []rune("<i class='inline'>"), Priority: i, IsOpen: true})
			markUpByPosition = append(markUpByPosition, MarkupNyPosition{Offset: entity.Offset + entity.Length, Tag: []rune("</i>"), Priority: i, IsOpen: false})
		case telebot.EntityBold:
			markUpByPosition = append(markUpByPosition, MarkupNyPosition{Offset: entity.Offset, Tag: []rune("<b class='inline'>"), Priority: i, IsOpen: true})
			markUpByPosition = append(markUpByPosition, MarkupNyPosition{Offset: entity.Offset + entity.Length, Tag: []rune("</b>"), Priority: i, IsOpen: false})
		case telebot.EntityURL:
			link := string(utf16.Decode(text[entity.Offset : entity.Offset+entity.Length]))
			if strings.Contains(link, "://") == false {
				link = "http://" + link
			}
			markUpByPosition = append(markUpByPosition, MarkupNyPosition{Offset: entity.Offset, Tag: []rune("<a target='_blank' href='" + link + "' class='inline c-link'>"), Priority: i, IsOpen: true})
			markUpByPosition = append(markUpByPosition, MarkupNyPosition{Offset: entity.Offset + entity.Length, Tag: []rune("</a>"), Priority: i, IsOpen: false})
		case telebot.EntityTextLink:
			markUpByPosition = append(markUpByPosition, MarkupNyPosition{Offset: entity.Offset, Tag: []rune("<a target='_blank' class='inline c-link' href='" + entity.URL + "'>"), Priority: i, IsOpen: true})
			markUpByPosition = append(markUpByPosition, MarkupNyPosition{Offset: entity.Offset + entity.Length, Tag: []rune("</a>"), Priority: i, IsOpen: false})
		case telebot.EntityMention:
			link := string(utf16.Decode(text[entity.Offset+1 : entity.Offset+entity.Length]))
			markUpByPosition = append(markUpByPosition, MarkupNyPosition{Offset: entity.Offset, Tag: []rune("<a target='_blank' href='https://t.me/" + link + "' class='inline c-link'>"), Priority: i, IsOpen: true})
			markUpByPosition = append(markUpByPosition, MarkupNyPosition{Offset: entity.Offset + entity.Length, Tag: []rune("</a>"), Priority: i, IsOpen: false})
		default:
			continue
		}
	}

	sort.Slice(markUpByPosition, func(i, j int) bool {
		a := markUpByPosition[i]
		b := markUpByPosition[j]

		// # If they are on the same place
		if a.Offset == b.Offset {
			// # If it is closing tag, then more prior must be first
			if a.IsOpen == false && b.IsOpen == false {
				return a.Priority < b.Priority
			}

			return a.Priority > b.Priority
		} else {
			return a.Offset > b.Offset
		}
	})

	for _, markup := range markUpByPosition {
		text = slices.Insert(text, markup.Offset, utf16.Encode(markup.Tag)...)
	}

	return string(utf16.Decode(text)), nil
}
