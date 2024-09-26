package domains

import "auth-service/models"

type ThirdPartyRepository interface {
	Create(user *models.AuthThirdParty) error
}
