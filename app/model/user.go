// package model

// import (
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// )

// type User struct {
//     ID        int       `json:"id"`
//     Username  string    `json:"username"`
//     Email     string    `json:"email"`
//     Role      string    `json:"role"`
//     CreatedAt time.Time `json:"created_at"`
// }

// type LoginRequest struct {
//     Username string `json:"username"`
//     Password string `json:"password"`
// }

// type LoginResponse struct {
//     User  User   `json:"user"`
//     Token string `json:"token"`
// }

// type JWTClaims struct {
//     UserID   int    `json:"user_id"`
//     Username string `json:"username"`
//     Role     string `json:"role"`
//     jwt.RegisteredClaims
// }

package model

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Username  string             `bson:"username" json:"username"`
    Email     string             `bson:"email" json:"email"`
    Role      string             `bson:"role" json:"role"`
    CreatedAt string          `bson:"created_at" json:"created_at"`
    PasswordHash string             `bson:"password_hash" json:"-"`
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    User  User   `json:"user"`
    Token string `json:"token"`
}

type JWTClaims struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

