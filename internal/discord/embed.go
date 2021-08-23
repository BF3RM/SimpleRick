package discord

//type EmbedOption func(embed *Embed)
//
//func WithTitle(title string) EmbedOption {
//	return func(embed *Embed) {
//		embed.Title = title
//	}
//}
//
//func WithDescription(description string) EmbedOption {
//	return func(embed *Embed) {
//		embed.Description = description
//	}
//}
//
//func WithURL(url string) EmbedOption {
//	return func(embed *Embed) {
//		embed.URL = url
//	}
//}
//
//func WithTimestamp() EmbedOption {
//	return func(embed *Embed) {
//		embed.Title = time.Now().Format(time.RFC3339)
//	}
//}
//
//func WithColor(color int) EmbedOption {
//	return func(embed *Embed) {
//		embed.Color = color
//	}
//}
//
//func WithProvider(name string, url string) EmbedOption {
//	return func(embed *Embed) {
//		embed.Provider = &EmbedProvider{
//			Name: name,
//			URL:  url,
//		}
//	}
//}
//
//func WithField(name string, value string, inline bool) EmbedOption {
//	return func(embed *Embed) {
//		embed.Fields = append(embed.Fields, &EmbedField{
//			Name:   name,
//			Value:  value,
//			Inline: inline,
//		})
//	}
//}
//
//func NewEmbed(opts ...EmbedOption) *Embed {
//	embed := &Embed{}
//
//	for _, opt := range opts {
//		opt(embed)
//	}
//
//	return embed
//}
//
