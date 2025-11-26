package utils

import (
	"crud_alumni/app/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret-key-panjang-minimal-32-char")


func GenerateToken(user model.User) (string, error) {
    claims := model.JWTClaims{
        UserID:   user.ID.Hex(), // ubah ObjectID ke string Hex
        Username: user.Username,
        Role:     user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// func GenerateToken(user model.User) (string, error) {
//     claims := model.JWTClaims{
//         UserID:   user.ID,
//         Username: user.Username,
//         Role:     user.Role,
//         RegisteredClaims: jwt.RegisteredClaims{
//             ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
//             IssuedAt:  jwt.NewNumericDate(time.Now()),
//         },
//     }

//     token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//     return token.SignedString(jwtSecret)
// }

func ValidateToken(tokenString string) (*model.JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*model.JWTClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, jwt.ErrInvalidKey
}
