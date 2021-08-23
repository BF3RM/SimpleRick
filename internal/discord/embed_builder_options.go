package discord

type EmbedFooterOption func(footer *EmbedFooter)

func WithFooterIcon(iconUrl string) EmbedFooterOption {
	return func(footer *EmbedFooter) {
		footer.IconUrl = iconUrl
	}
}

type EmbedImageOption func(image *EmbedImage)

func WithImageSize(width, height int) EmbedImageOption {
	return func(image *EmbedImage) {
		image.Width = width
		image.Height = height
	}
}

type EmbedAuthorOption func(author *EmbedAuthor)

func WithAuthorUrl(url string) EmbedAuthorOption {
	return func(author *EmbedAuthor) {
		author.URL = url
	}
}

func WithAuthorIcon(iconUrl string) EmbedAuthorOption {
	return func(author *EmbedAuthor) {
		author.IconUrl = iconUrl
	}
}
