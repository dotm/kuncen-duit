package userSignIn

import (
	"context"
	"encoding/json"
	"fmt"
	"kuncenduit-backend/shared/commandresponse"
	"kuncenduit-backend/shared/createerror"
	"kuncenduit-backend/shared/jwttoken"
	"kuncenduit-backend/shared/lazylogger"
	"kuncenduit-backend/shared/password"
	"kuncenduit-backend/shared/postgredb"
	"kuncenduit-backend/shared/stringmasker"
	"kuncenduit-backend/src/user"
	"time"

	"github.com/jackc/pgx/v4"
)

/*
Commands represent input from client through API requests.
Addition, change, or removal of struct fields might cause version increment
*/
type CommandV1 struct {
	Version  string `json:"version"` //should follow the struct name suffix
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (x CommandV1) createLoggableString() (string, error) {
	//strip any sensitive information.
	//strip any fields that are too large to be printed (e.g. image blob).
	loggableCommand := CommandV1{
		Version:  x.Version,
		Email:    stringmasker.Email(x.Email),
		Password: stringmasker.Password(x.Password),
	}
	byteSlice, err := json.Marshal(loggableCommand)
	if err != nil {
		return "", err
	} else {
		return string(byteSlice), nil
	}
}

type CommandV1DataResponse struct {
	SignedJwtToken string
	JwtExpiration  time.Time
}

/*
Addition, change, or removal of validation might cause version increment
*/
func CommandV1Handler(ctx context.Context, logger *lazylogger.Instance, command CommandV1) (CommandV1DataResponse, *commandresponse.ErrorObj) {
	//don't mutate this. emptyResponse should be used when returning error.
	emptyResponse := CommandV1DataResponse{}
	//log the command
	loggableCommand, err := command.createLoggableString()
	if err != nil {
		err = fmt.Errorf("error creating loggable string: %v", err)
		logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}
	logger.EnqueueCommandLog(loggableCommand, true)

	/* Validations
	Validations from auth, write model,
	or domain model's business logic (from projections or from events replay).
	*/
	conn, err := postgredb.AcquireConnectionToMainDatabase(ctx)
	if err != nil {
		err = fmt.Errorf("error connecting to database: %v", err)
		logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}
	defer conn.Close(ctx)

	var existingUser user.Domain
	var id string
	var name string
	var email string
	var hashedPassword string
	err = conn.QueryRow(
		ctx,
		"SELECT id, name, email, hashed_password FROM public.users where email=$1",
		command.Email).
		Scan(&id, &name, &email, &hashedPassword)
	if err == pgx.ErrNoRows {
		return emptyResponse, createerror.UserCredentialIncorrect()
	}
	if err != nil {
		err = fmt.Errorf("error executing query: %v", err)
		logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}
	existingUser = user.Domain{
		Id:             id,
		Name:           &name,
		Email:          email,
		HashedPassword: &hashedPassword,
	}

	passwordCorrect := password.MatchPasswordToHash(command.Password, *existingUser.HashedPassword)
	if !passwordCorrect {
		return emptyResponse, createerror.UserCredentialIncorrect()
	}

	/* Business Logic
	Perform business logic preferably through domain model's methods.
	*/
	jwtExpiration := time.Now().Add(time.Hour * 24) //1 day
	jwtClaims := map[string]string{"user_id": existingUser.Id}
	signedJwtToken, err := jwttoken.BuildAndSign(jwtClaims, jwtExpiration)
	if err != nil {
		logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}

	/* Persisting Data
	Persist event to event store.
	If write model is used, also persist write model with atomic transaction.
	*/

	//You can send the event id back to the requester
	//so that they can periodically check the status of the event.
	return CommandV1DataResponse{
		SignedJwtToken: signedJwtToken,
		JwtExpiration:  jwtExpiration,
	}, nil
}
