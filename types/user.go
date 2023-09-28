package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinFirstNameLen = 2
	MinLastNameLen  = 2
	MinPasswordLen  = 7
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func CreateUserFromParams(params *CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), 12)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}

func (params *CreateUserParams) Validate() []string {
	errors := []string{}
	if len(params.FirstName) < MinFirstNameLen {
		errors = append(errors, fmt.Sprintf("firstname length should be at least %d characters", MinFirstNameLen))
	}
	if len(params.LastName) < MinLastNameLen {
		errors = append(errors, fmt.Sprintf("lastname length should be at least %d characters", MinLastNameLen))
	}
	if !isValidEmail(params.Email) {
		errors = append(errors, fmt.Sprintf("Email is invalid"))
	}
	if len(params.Password) < MinPasswordLen {
		errors = append(errors, fmt.Sprintf("Password length should be at least %d characters", MinPasswordLen))
	}
	return errors
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func (params *UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}
	if len(params.FirstName) > 0 {
		m["firstName"] = params.FirstName
	}
	if len(params.LastName) > 0 {
		m["lastName"] = params.LastName
	}
	return m
}
