package user

import (
	"encoding/json"
	"errors"
	enUser "abishar/internal/entity/user"
	"abishar/internal/pkg/encrypt"
)

func (uc *Usecase) SaveSession(user *enUser.User) (string, error) {
  uniqueKey, err := encrypt.GenerateUUID()
	if err != nil {
		return "", errors.New("[SaveSession] failed to generate UUID") 
  }

  session := enUser.Session{
    ID: user.ID,
    Username: user.Username,
    UniqueKey: uniqueKey,
    IsAdmin: user.IsAdmin,
  }

  token, err := uc.userRepo.GenerateSessionToken(session)
  if err != nil {
    return "", errors.New("[SaveSession] failed to generate token")
  }

  sessionData, _ := json.Marshal(enUser.SessionData{
    ID: user.ID,
    Token: token,
  })

  if err = uc.userRepo.SaveSession(session, int(enUser.TokenTTL.Seconds()), sessionData); err != nil {
    return "", errors.New("[SaveSession] failed to save session")
  }

  return token, nil
}

func (uc *Usecase) GetUserSession(sess enUser.Session) *enUser.SessionData {
  result := uc.userRepo.GetSession(sess)

  sessionData := &enUser.SessionData{}
  resultByte, ok := result.Value.([]byte)
  if !ok {
    return nil
  }
  if err := json.Unmarshal(resultByte, sessionData); err != nil {
    return nil 
  }

  return sessionData
} 
