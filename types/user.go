package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost         = 12
	minFirstNameLenght = 2
	minLastNameLenght  = 2
	minPasswordLenght  = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (c *CreateUserParams) Validation() []string {
	var errors []string
	if len(c.FirstName) < minFirstNameLenght {
		errors = append(errors, fmt.Sprintf("FirstName length must be greater then %d.", minFirstNameLenght))
	}
	if len(c.LastName) < minLastNameLenght {
		errors = append(errors, fmt.Sprintf("LastName length must be greater then %d.", minLastNameLenght))
	}
	if len(c.Password) < minPasswordLenght {
		errors = append(errors, fmt.Sprintf("Password length must be greater then %d.", minPasswordLenght))
	}
	if !IsValidEmail(c.Email) {
		errors = append(errors, "Email is InValid!")
	}
	return errors
}

func IsValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encryptedPass),
	}, nil
}
