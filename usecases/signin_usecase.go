package usecases

import "auth-service/domains"

type signInUsecase struct {
}

func NewSignInUsecase() domains.SignInUsecase {
	return &signInUsecase{}
}
