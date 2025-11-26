// package service

// import (
// 	"crud_alumni/app/model"
// 	"crud_alumni/app/repository"
// 	"crud_alumni/utils"
// 	"errors"
// 	"fmt"

// 	"golang.org/x/crypto/bcrypt"
// )

// func Login(req model.LoginRequest) (*model.LoginResponse, error) {
//     fmt.Println("=== DEBUG LOGIN ===")
//     fmt.Println("Input Username:", req.Username)
//     fmt.Println("Input Password:", req.Password)
//     hashInput, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// fmt.Println("Password Input Hashed (baru):", string(hashInput))

//     user, passwordHash, err := repository.FindUserByUsernameOrEmail(req.Username)
//     if err != nil {
//         fmt.Println("DB Error:", err)
//         return nil, errors.New("username atau password salah")
//     }

//     fmt.Println("User from DB:", user.Username, user.Email, user.Role)
//     fmt.Println("Password Hash from DB:", passwordHash)

//     ok := utils.CheckPassword(req.Password, passwordHash)
//     fmt.Println("Password Match:", ok)

//     if !ok {
//         return nil, errors.New("username atau password salah")
//     }

//     token, err := utils.GenerateToken(*user)
//     if err != nil {
//         return nil, errors.New("gagal generate token")
//     }

//     return &model.LoginResponse{
//         User:  *user,
//         Token: token,
//     }, nil
// }

package service

import (
	"crud_alumni/app/model"
	"crud_alumni/app/repository"
	"crud_alumni/utils"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// LoginMongo - versi login untuk MongoDB dengan debug hash
func Login(req model.LoginRequest) (*model.LoginResponse, error) {
	fmt.Println("=== DEBUG LOGIN ===")
	fmt.Println("Input Username:", req.Username)
	fmt.Println("Input Password:", req.Password)

	// 1. Ambil user + hash password dari MongoDB
	user, passwordHashDB, err := repository.FindUserByUsernameOrEmail(req.Username)
	if err != nil {
		fmt.Println("User tidak ditemukan atau DB error:", err)
		return nil, errors.New("username atau password salah")
	}

	fmt.Println("User dari DB:", user.Username, user.Email, user.Role)
	fmt.Println("Password hash DB:", passwordHashDB)

	// 2. Compare password input dengan hash DB
	err = bcrypt.CompareHashAndPassword([]byte(passwordHashDB), []byte(req.Password))
	if err != nil {
		fmt.Println("Password tidak cocok ❌")
		return nil, errors.New("username atau password salah")
	}

	fmt.Println("Password cocok ✅")

	// 3. Generate token JWT
	token, err := utils.GenerateToken(*user)
	if err != nil {
		fmt.Println("Gagal generate token:", err)
		return nil, errors.New("gagal generate token")
	}

	// 4. Return response
	return &model.LoginResponse{
		User:  *user,
		Token: token,
	}, nil
}


