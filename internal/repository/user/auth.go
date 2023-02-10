package user

import (
	enUser "abishar/internal/entity/user"
	"abishar/internal/pkg/redigo"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func (r *Repository) GenerateSessionToken(sess enUser.Session) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		enUser.JWTFieldID:         sess.ID,
		enUser.JWTFieldUsername:   sess.Username,
		enUser.JWTFieldUniequeKey: sess.UniqueKey,
		enUser.JWTFieldIsAdmin:    sess.IsAdmin,
		"nbf":                     time.Now().Unix(),
		"exp":                     time.Now().Add(enUser.TokenTTL).Unix(),
	})

	return token.SignedString([]byte(r.jwt.Secret))
}

func (r *Repository) SaveSession(sess enUser.Session, expireTime int, data []byte) error {
	return r.redis.Setex(fmt.Sprintf("session:%d:%s", sess.ID, sess.UniqueKey), expireTime, data)
}

func (r *Repository) RemoveSession(sess enUser.Session) error {
	return r.redis.Del(fmt.Sprintf("session:%d:%s", sess.ID, sess.UniqueKey))
}

func (r *Repository) GetSession(sess enUser.Session) *redigo.Result {
	return r.redis.Get(fmt.Sprintf("session:%d:%s", sess.ID, sess.UniqueKey))
}
