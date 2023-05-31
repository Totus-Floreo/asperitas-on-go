package helpers

import "github.com/Totus-Floreo/asperitas-on-go/internal/model"

func HTTPError(err error) string {
	switch err {
	case model.ErrPostNotFound:
		return model.ErrPostNotFoundHTTP.Error()
	case model.ErrInvalidPostID:
		return model.ErrPostInvalidHTTP.Error()
	case model.ErrInvalidCredentials:
		return model.ErrInvalidCredentialsHTTP.Error()
	case model.ErrUnAuthorized:
		return model.ErrUnAuthorizedHTTP.Error()
	case model.ErrInvalidCommentID:
		return model.ErrCommentInvalidHTTP.Error()
	}
	return err.Error()
}
