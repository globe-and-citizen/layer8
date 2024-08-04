package code

import "globe-and-citizen/layer8/server/resource_server/models"

type Generator interface {
	GenerateCode(user *models.User, emailAddress string) (string, error)
}
