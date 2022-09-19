package createerror

import "kuncenduit-backend/shared/commandresponse"

func UserNotFound(err error) *commandresponse.ErrorObj {
	return Response(
		"user/not-found",
		err,
		map[string]bool{
			commandresponse.SendErrorToLog:      true,
			commandresponse.SendErrorToDevEmail: true,
		},
	)
}

const UserCredentialIncorrectErrorCode = "user/credential-incorrect"

func UserCredentialIncorrect() *commandresponse.ErrorObj {
	return Response(
		UserCredentialIncorrectErrorCode,
		nil,
		map[string]bool{},
	)
}

const UserCredentialExpiredErrorCode = "user/credential-expired"

func UserCredentialExpired() *commandresponse.ErrorObj {
	return Response(
		UserCredentialExpiredErrorCode,
		nil,
		map[string]bool{},
	)
}
