package discord

import "time"

type EmbedBuilder struct {
	embed Embed
}

func NewEmbedBuilder() *EmbedBuilder {
	return &EmbedBuilder{}
}

func (e *EmbedBuilder) SetTitle(title string) *EmbedBuilder {
	e.embed.Title = title
	return e
}

func (e *EmbedBuilder) SetDescription(description string) *EmbedBuilder {
	e.embed.Description = description
	return e
}

func (e *EmbedBuilder) SetURL(url string) *EmbedBuilder {
	e.embed.URL = url
	return e
}

func (e *EmbedBuilder) AddTimestamp() *EmbedBuilder {
	e.embed.Timestamp = time.Now().Format(time.RFC3339)
	return e
}

func (e *EmbedBuilder) SetColor(color int) *EmbedBuilder {
	e.embed.Color = color
	return e
}

func (e *EmbedBuilder) SetAuthor(name string, opts ...EmbedAuthorOption) *EmbedBuilder {
	author := &EmbedAuthor{
		Name: name,
	}

	for _, opt := range opts {
		opt(author)
	}

	e.embed.Author = author

	return e
}

func (e *EmbedBuilder) SetThumbnail(url string, opts ...EmbedImageOption) *EmbedBuilder {
	thumbnail := &EmbedImage{
		URL: url,
	}

	for _, opt := range opts {
		opt(thumbnail)
	}

	e.embed.Thumbnail = thumbnail

	return e
}

func (e *EmbedBuilder) SetImage(url string, opts ...EmbedImageOption) *EmbedBuilder {
	image := &EmbedImage{
		URL: url,
	}

	for _, opt := range opts {
		opt(image)
	}

	e.embed.Image = image

	return e
}

func (e *EmbedBuilder) SetVideo(url string, opts ...EmbedImageOption) *EmbedBuilder {
	video := &EmbedImage{
		URL: url,
	}

	for _, opt := range opts {
		opt(video)
	}

	e.embed.Video = video

	return e
}

type EmbedFieldOption func(field *EmbedField)

func WithFieldInline() EmbedFieldOption {
	return func(field *EmbedField) {
		field.Inline = true
	}
}

func (e *EmbedBuilder) AddField(name, value string, opts ...EmbedFieldOption) *EmbedBuilder {
	field := &EmbedField{
		Name:  name,
		Value: value,
	}

	for _, opt := range opts {
		opt(field)
	}

	e.embed.Fields = append(e.embed.Fields, field)

	return e
}

func (e *EmbedBuilder) SetFooter(text string, opts ...EmbedFooterOption) *EmbedBuilder {
	e.embed.Footer = &EmbedFooter{
		Text: text,
	}

	for _, opt := range opts {
		opt(e.embed.Footer)
	}

	return e
}

func (e EmbedBuilder) Build() Embed {
	return e.embed
}
