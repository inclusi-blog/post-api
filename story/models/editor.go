package models

import "github.com/mitchellh/mapstructure"

type ElementType string

const (
	Paragraph ElementType = "paragraph"
	Header    ElementType = "header"
	Table     ElementType = "table"
	List      ElementType = "list"
	Quote     ElementType = "quote"
	CheckList ElementType = "checklist"
	Warning   ElementType = "warning"
	Code      ElementType = "code"
	LinkTool  ElementType = "linkTool"
	Image     ElementType = "image"
	RawHTML   ElementType = "raw"
	Separator ElementType = "delimiter"
)

func (elementType ElementType) IsEqual(otherElementType ElementType) bool {
	return elementType == otherElementType
}

type Block struct {
	ID   string                 `json:"id"`
	Type ElementType            `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type Editor struct {
	Time    int64   `json:"time"`
	Blocks  []Block `json:"blocks"`
	Version string  `json:"version"`
}

type ImageElement struct {
	File struct {
		Url string `json:"url"`
	} `json:"file"`
	Caption        string `json:"caption"`
	WithBorder     bool   `json:"withBorder"`
	Stretched      bool   `json:"stretched"`
	WithBackground bool   `json:"withBackground"`
}

type ParagraphElement struct {
	Text string `json:"text"`
}

func (block Block) GetText() (string, error) {
	if block.Type.IsEqual(Paragraph) {
		var paragraph ParagraphElement
		err := mapstructure.Decode(block.Data, &paragraph)
		if err != nil {
			return "", err
		}
		return paragraph.Text, nil
	}
	var header HeaderElement
	err := mapstructure.Decode(block.Data, &header)
	if err != nil {
		return "", err
	}
	return header.Text, nil
}

type HeaderElement struct {
	Text  string `json:"text"`
	Level int    `json:"level"`
}

func (e *Editor) WithImageElement(element ImageElement) {
	e.Blocks = append(e.Blocks, Block{Type: Image})
}

func (e *Editor) WithParagraphElement(element Block) {
	e.Blocks = append(e.Blocks, Block{Type: Paragraph})
}

func (e *Editor) WithListElement(element Block) {
	e.Blocks = append(e.Blocks, Block{Type: List})
}
