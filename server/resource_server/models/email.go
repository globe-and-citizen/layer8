package models

type Email struct {
	From    string
	To      string
	Subject string
	Content VerificationEmailContent
}

type VerificationEmailContent struct {
	Username string
	Code     string
}
