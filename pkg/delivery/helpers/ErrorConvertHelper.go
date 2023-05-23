package helpers

import "github.com/Totus-Floreo/asperitas-on-go/pkg/model"

func HTTPError(err error) string {
	switch err {
	case model.ErrPostNotFound:
		return model.ErrPostInvalidHTTP.Error()
	case model.ErrPostNotFound:
		return model.ErrPostInvalidHTTP.Error()
	case model.ErrInvalidCredentials:
		return model.ErrInvalidCredentialsHTTP.Error()
	case model.ErrUnAuthorized:
		return model.ErrUnAuthorizedHTTP.Error()
	}
	return err.Error()
}
