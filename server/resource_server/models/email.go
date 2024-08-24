package models

type Email struct {
	SenderAddress        string
	RecipientAddress     string
	RecipientDisplayName string
	Subject              string
	Content              string
}
