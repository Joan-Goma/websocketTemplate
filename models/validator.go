package models

import (
	"errors"
	engine "github.com/JoanGTSQ/api"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
)

type userValidator struct {
	UserDB
	hmac       engine.HMAC
	emailRegex *regexp.Regexp
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

var uuserValidator userValidator

func newUserValidator(udb UserDB, hmac engine.HMAC) *userValidator {
	uuserValidator.UserDB = udb
	uuserValidator.hmac = hmac
	uuserValidator.emailRegex = regexp.MustCompile(`^[a-z0-9._%+\]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

func bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}

	pwByte := []byte(user.Password + userPwPPepper)

	// Generate hash from password and use hash cost
	hashedBytes, err := bcrypt.GenerateFromPassword(pwByte, hashCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return engine.ERR_PSSWD_TOO_SHORT
	}
	return nil
}

func passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return engine.ERR_PSSWD_REQUIRED
	}
	return nil
}

func passwordRequired(user *User) error {
	if user.Password == "" {
		engine.Warning.Println("No se ha encontrado ninguna psswd")
		return engine.ERR_PSSWD_REQUIRED
	}
	return nil
}

func hmacRemember(user *User) error {
	if user.Remember == "" {
		return errors.New("could not find remember")
	}
	user.RememberHash = uuserValidator.hmac.Hash(user.Remember)
	return nil
}

func defaultify(user *User) error {
	if user.Remember != "" {
		return errors.New("could not find remember")
	}

	token, err := engine.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func defaultifyCreation(user *User) error {
	user.RoleID = engine.ROLE_USER
	user.Banned = false
	return nil
}

func rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := engine.NBytes(user.Remember)
	if err != nil {
		return err

	}
	if n < 32 {
		return engine.ERR_REMMEMBER_TOO_SHOT
	}
	return nil
}
func rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return engine.ERR_REMMEMBER_REQUIRED
	}
	return nil
}

func normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)

	return nil
}

func requireEmail(user *User) error {
	if user.Email == "" {
		return engine.ERR_MAIL_REQUIRED
	}
	return nil
}

func emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uuserValidator.emailRegex.MatchString(user.Email) {
		return engine.ERR_MAIL_IS_N0T_VALID
	}
	return nil
}

func emailsIsAvail(user *User) error {
	err := user.ByEmail()

	switch err {
	case engine.ERR_NOT_FOUND:
		return nil
	case nil:
		return engine.ERR_MAIL_IS_TAKEN
	default:
		return engine.ERR_MAIL_NOT_EXIST
	}
}

func (user *User) NewUserValidation() error {
	if err := runUserValFuncs(user,
		passwordRequired,
		passwordMinLength,
		bcryptPassword,
		passwordHashRequired,
		defaultify,
		defaultifyCreation,
		hmacRemember,
		rememberHashRequired,
		normalizeEmail,
		requireEmail,
		emailFormat,
		emailsIsAvail); err != nil {
		return err
	}

	return nil
}
func (user *User) ExistentUserValidation() error {
	if err := runUserValFuncs(user,
		passwordRequired,
		passwordMinLength,
		bcryptPassword,
		passwordHashRequired,
		defaultify,
		hmacRemember,
		rememberHashRequired,
		normalizeEmail,
		requireEmail,
		emailFormat,
	); err != nil {
		return err
	}

	return nil
}
