package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	enUser "abishar/internal/entity/user"

	"abishar/internal/pkg/encrypt"
	"abishar/internal/pkg/redigo"
)

type (
	userRepository interface {
		InsertUser(ctx context.Context, form enUser.RegisterForm) (*enUser.User, error)
		CheckUsername(ctx context.Context, username string) (int64, error)
		GenerateSessionToken(sess enUser.Session) (string, error)
		SaveSession(sess enUser.Session, expireTime int, data []byte) error
		GetByUsername(ctx context.Context, username string) (*enUser.User, error)
    RemoveSession(sess enUser.Session) error
    GetSession(sess enUser.Session) *redigo.Result
    GetUserWallet(ctx context.Context, username string) (*enUser.UserWallet, error)
    AddWallet(ctx context.Context, amount int64, userID int64) error
	}
)

type Usecase struct {
	userRepo userRepository
}

func NewUsecase(
	userRepo userRepository,
) *Usecase {
	return &Usecase{
		userRepo: userRepo,
	}
}

func (uc *Usecase) RegisterUser(ctx context.Context, form enUser.RegisterForm) (*enUser.RegisterResponse, error) {

	var errs []error

	// Field rule validation
	if len(form.Username) < 6 {
		errs = append(errs, errors.New("Username required to be more than 6 chars"))
	}

	if len(form.Password) < 6 {
		errs = append(errs, errors.New("Password required to be more than 6 chars"))
	}

	// Check username exist
	userID, err := uc.userRepo.CheckUsername(ctx, form.Username)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to check username availabilty. err %v", err.Error()))
	}

	if userID != 0 {
		errs = append(errs, errors.New("Username already existed"))
	}

  if len(errs) != 0 {
    return nil, errs[0]
  }

	form.Salt, err = encrypt.GenerateUUID()
	if err != nil {
		log.Printf("[RegisterUser] failed to generate UUID. Err: %v", err)
		return nil, err
	}

	form.Password = encrypt.EncodeSHA1(form.Password + form.Salt)

	user, err := uc.userRepo.InsertUser(ctx, form)
	if err != nil {
		return nil, err
	}

	token, err := uc.SaveSession(user)

	return &enUser.RegisterResponse{
		ID:    user.ID,
		Token: token,
	}, nil
}

func (uc *Usecase) Login(ctx context.Context, form enUser.LoginRequest) (*enUser.RegisterResponse, error) {

  var errs []error

	// Field rule validation
	if len(form.Username) < 6 {
		errs = append(errs, errors.New("Username required to be more than 6 chars"))
	}

	if len(form.Password) < 6 {
		errs = append(errs, errors.New("Password required to be more than 6 chars"))
	}

  user, err := uc.userRepo.GetByUsername(ctx, form.Username)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Failed to get user. err %v", err.Error()))
  }

  if user.ID == 0 {
    errs = append(errs, errors.New("Username does not exist"))
  } else if user.Password != encrypt.EncodeSHA1(form.Password+user.Salt) {
    errs = append(errs, errors.New("Wrong Password"))
  }

  if len(errs) != 0 {
    return nil, errs[0]
  }

  token, err := uc.SaveSession(user)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Failed to save session. err %v", err.Error()))
  }

  return &enUser.RegisterResponse{
    ID: user.ID,
    Token: token,
  }, nil

}

func (uc *Usecase) Logout(sess enUser.Session) error {

  err := uc.userRepo.RemoveSession(sess)
  if err != nil {
    return errors.New(fmt.Sprintf("[Logout] Failed to log out. err: %v", err.Error()))
  }

  return nil
}

func (uc *Usecase) GetUserWallet(ctx context.Context, username string) (*enUser.UserWallet, error) {
  user, err := uc.userRepo.GetUserWallet(ctx, username)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("[GetUserWallet] Failed to get user. err: %v", err.Error()))
  }

  return user, nil
}

func (uc *Usecase) AddWallet(ctx context.Context, amount int64, userID int64) error {
  err := uc.userRepo.AddWallet(ctx, amount, userID)
  if err != nil {
    return errors.New(fmt.Sprintf("[AddWallet] Failed to add wallet. err: %v", err.Error()))
  }

  return nil
}
