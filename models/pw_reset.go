package models

import (
	engine "github.com/JoanGTSQ/api"
	"github.com/jinzhu/gorm"
)

type pwReset struct {
	gorm.Model
	UserID    uint   `gorm:"not null"`
	Token     string `gorm:"-"`
	TokenHash string `gorm:"not null;unique_index"`
}

type pwResetDB interface {
	ByToken(token string) (*pwReset, error)
	ByUserID(userID uint) (*pwReset, error)

	Create(pwr *pwReset) error
	Delete(id uint) error
}
type pwResetValidator struct {
	pwResetDB
	hmac engine.HMAC
}
type pwResetGorm struct {
	db *gorm.DB
}

func (pwrv *pwResetValidator) ByToken(token string) (*pwReset, error) {
	pwr := pwReset{Token: token}
	err := runPwResetValFns(&pwr, pwrv.hmacToken)
	if err != nil {
		return nil, err
	}
	return pwrv.pwResetDB.ByToken(pwr.Token)
}

func (pwrv *pwResetValidator) Create(pwr *pwReset) error {
	err := runPwResetValFns(pwr,
		pwrv.requireUserID,
		pwrv.setTokenIfUnset,
		pwrv.hmacToken)
	if err != nil {
		return err
	}
	return pwrv.pwResetDB.Create(pwr)
}

func (pwrv *pwResetValidator) Delete(id uint) error {
	if id <= 0 {
		return engine.ERR_ID_INVALID
	}
	return pwrv.pwResetDB.Delete(id)
}

func (pwrg *pwResetGorm) ByToken(tokenHash string) (*pwReset, error) {
	var pwr pwReset
	err := first(pwrg.db.Where("token_hash = ?", tokenHash), &pwr)
	if err != nil {
		return nil, err
	}
	return &pwr, nil
}
func (pwrg *pwResetGorm) ByUserID(userID uint) (*pwReset, error) {
	var pwr pwReset
	err := first(pwrg.db.Where("user_id = ?", userID), &pwr)
	if err != nil {
		return nil, err
	}
	return &pwr, nil
}
func (pwrg *pwResetGorm) Create(pwr *pwReset) error {
	return pwrg.db.Create(pwr).Error
}

func (pwrg *pwResetGorm) Delete(id uint) error {
	pwr := pwReset{Model: gorm.Model{ID: id}}
	return pwrg.db.Delete(&pwr).Error
}

func (pwrv *pwResetValidator) requireUserID(pwr *pwReset) error {
	if pwr.UserID <= 0 {
		return engine.ERR_USERID_REQUIRED
	}
	return nil
}

func (pwrv *pwResetValidator) setTokenIfUnset(pwr *pwReset) error {
	if pwr.Token != "" {
		return nil
	}
	token, err := engine.RememberToken()
	if err != nil {
		return err
	}
	pwr.Token = token
	return nil
}

func (pwrv *pwResetValidator) hmacToken(pwr *pwReset) error {
	if pwr.Token == "" {
		return nil
	}
	pwr.TokenHash = pwrv.hmac.Hash(pwr.Token)
	return nil
}

func newPwResetValidator(db pwResetDB, hmac engine.HMAC) *pwResetValidator {
	return &pwResetValidator{
		db,
		hmac,
	}
}

type pwResetValFn func(reset *pwReset) error

func runPwResetValFns(pwr *pwReset, fns ...pwResetValFn) error {
	for _, fn := range fns {
		if err := fn(pwr); err != nil {
			return err
		}
	}
	return nil
}
