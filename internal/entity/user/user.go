package user

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	JWTFieldID         = "trewmi"
	JWTFieldUsername   = "tmadqj"
	JWTFieldUniequeKey = "yim28m"
	JWTFieldIsAdmin    = "jeqhuh"

	TokenTTL = 24 * time.Hour

	ErrorJWTHasExpired       = "jwt has expired"
	ErrorInvalidJWTClaimID   = "invalid jwt claim id"
	ErrorInvalidJWTClaimType = "invalid jwt claim type"

	SessionContextKey = "authjwt"
)

type TokenClaim struct {
	ID        int64  `json:"trewmi"`
	Username  string `json:"tmadqj"`
	UniqueKey string `json:"yim28m"`
	IsAdmin   bool   `json:"jeqhuh"`

	jwt.StandardClaims
}

func (t *TokenClaim) Valid() error {

	currentTime := time.Now()
	expiresAt := time.Unix(t.ExpiresAt, 0)
	if expiresAt.Sub(currentTime) < 0 {
		// If token has expired
		return errors.New(ErrorJWTHasExpired)
	} else if t.ID <= 0 {
		// Invalid ID
		return errors.New(ErrorInvalidJWTClaimID)
	}

	return nil
}

type RegisterForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"-"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID    int64  `json:"userID"`
	Token string `json:"token"`
}

type Session struct {
	ID        int64
	Username  string
	UniqueKey string
	IsAdmin   bool
}

type SessionData struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
}

type User struct {
	ID       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"-" db:"password"`
	IsAdmin  bool   `json:"isAdmin" db:"is_admin"`
	Wallet   int64  `json:"wallet" db:"wallet"`
	Salt     string `json:"-" db:"salt"`
}

type WalletRequest struct {
	Amount int64 `json:"amount"`
}

type UserWallet struct {
	ID       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Wallet   int64  `json:"wallet" db:"wallet"`
}
