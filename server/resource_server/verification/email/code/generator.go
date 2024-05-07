package code

type Generator interface {
	GenerateCode(emailAddress string) string
}
