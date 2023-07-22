package authentication_user

import (
	"context"
	"errors"
	"time"

	"github.com/brutalzinn/boberto-modpack-api/common"
	config "github.com/brutalzinn/boberto-modpack-api/configs"
	"github.com/brutalzinn/boberto-modpack-api/infra/services/authentication/user/models"
	user_database "github.com/brutalzinn/boberto-modpack-api/repository/database/user"
	entities_user "github.com/brutalzinn/boberto-modpack-api/repository/database/user/entities"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(userId string) (string, error) {
	cfg := config.GetConfig()
	secretKey := config.GetJWTSecret()
	expirationTime := time.Now().Add(time.Duration(cfg.Authentication.Expiration) * time.Second)
	claims := models.Claims{
		ID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	return tokenString, err
}

func VerifyJWT(tokenJWT string) (models.Claims, error) {
	secretKey := config.GetJWTSecret()
	claims := models.Claims{}
	_, err := jwt.ParseWithClaims(tokenJWT, &claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return claims, errors.New("No valid token provided")
	}
	return claims, err
}

func GetCurrentUser(ctx context.Context) (user *entities_user.User, err error) {
	user_id, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, errors.New("No authorized to use this route")
	}
	user, err = user_database.Get(user_id)
	if err != nil {
		return nil, errors.New("No authorized to use this route")
	}
	return user, nil
}

func Authentication(email string, password string) (userId string, err error) {
	user, err := user_database.FindByEmail(email)
	if err != nil {
		return "", errors.New("Invalid user")
	}
	validPassword := common.BcryptCheckHash(password, user.Password)
	if !validPassword {
		return "", errors.New("Invalid user")
	}
	return user.ID, nil
}
