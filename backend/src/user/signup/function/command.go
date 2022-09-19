package userSignUp

import (
	"context"
	"encoding/json"
	"fmt"
	"kuncenduit-backend/shared/commandresponse"
	"kuncenduit-backend/shared/createerror"
	"kuncenduit-backend/shared/lazylogger"
	"kuncenduit-backend/shared/password"
	"kuncenduit-backend/shared/postgredb"
	"kuncenduit-backend/shared/stringmasker"
	"kuncenduit-backend/src/user"

	"github.com/google/uuid"
)

/*
Commands represent input from client through API requests.
Addition, change, or removal of struct fields might cause version increment
*/
type CommandV1 struct {
	Version  string `json:"version"` //should follow the struct name suffix
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (x CommandV1) createLoggableString() (string, error) {
	//strip any sensitive information.
	//strip any fields that are too large to be printed (e.g. image blob).
	loggableCommand := CommandV1{
		Version:  x.Version,
		Name:     stringmasker.Name(x.Name),
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

type CommandV1DataResponse = user.Domain

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
	//validate existing email should not register twice after moving to event sourcing ~kodok

	/* Business Logic
	Perform business logic preferably through domain model's methods.
	*/
	hashedPassword, err := password.Hash(command.Password)
	if err != nil {
		err = fmt.Errorf("error hashing password: %v", err)
		logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}

	newUser := user.Domain{
		Id:             uuid.NewString(),
		Email:          command.Email,
		Name:           &command.Name,
		HashedPassword: &hashedPassword,
	}

	/* Persisting Data
	Persist event to event store.
	If write model is used, also persist write model with atomic transaction.
	*/
	conn, err := postgredb.AcquireConnectionToMainDatabase(ctx)
	if err != nil {
		err = fmt.Errorf("error connecting to database: %v", err)
		logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}
	defer conn.Close(ctx)

	sqlStr := fmt.Sprintf(`
		INSERT INTO public.users
		(id, "name", email, "hashed_password")
		VALUES('%v', '%v', '%v', '%v');
	`, newUser.Id, *newUser.Name, newUser.Email, *newUser.HashedPassword)

	_, err = conn.Exec(ctx, sqlStr)
	if err != nil {
		err = fmt.Errorf("error executing query: %v", err)
		logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}

	//You can send the event id back to the requester
	//so that they can periodically check the status of the event.
	return newUser, nil
}
