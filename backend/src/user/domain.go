package user

//for use in command only. query doesn't need domain.
type Domain struct {
	Id    string `json:"id"`
	Email string `json:"email"`

	//any fields other than the necessary ones should be treated as nullable
	//so that new addition to the domain field can be treated as nullable too
	//with explicitly adding default value to legacy domain in persistence
	Name           *string `json:"name"`
	HashedPassword *string `json:"hashed_password"`

	// Nonce string `json:"nonce"` //used for concurrency control in NoSQL databases
}
