package discord

type EmbedImage struct {
	URL    string `json:"url,omitempty"`
	Height int    `json:"height,omitempty"`
	Width  int    `json:"width,omitempty"`
}

type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type EmbedAuthor struct {
	Name    string `json:"name,omitempty"`
	URL     string `json:"url,omitempty"`
	IconUrl string `json:"icon_url,omitempty"`
}

type EmbedFooter struct {
	Text    string `json:"text"`
	IconUrl string `json:"icon_url,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type Embed struct {
	Title       string         `json:"title,omitempty"`
	Description string         `json:"description,omitempty"`
	URL         string         `json:"url,omitempty"`
	Timestamp   string         `json:"timestamp,omitempty"`
	Color       int            `json:"color,omitempty"`
	Footer      *EmbedFooter   `json:"footer,omitempty"`
	Image       *EmbedImage    `json:"image,omitempty"`
	Thumbnail   *EmbedImage    `json:"thumbnail,omitempty"`
	Video       *EmbedImage    `json:"video,omitempty"`
	Provider    *EmbedProvider `json:"provider,omitempty"`
	Author      *EmbedAuthor   `json:"author,omitempty"`
	Fields      []*EmbedField  `json:"fields,omitempty"`
}

type WebhookPayload struct {
	Embeds []Embed `json:"embeds"`
}
