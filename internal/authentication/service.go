package authentication

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/johan-ag/wishlist/internal/users"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(registerDTO users.User) (bool, error)
	Login(loginDTO users.User) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

func NewAuthService(usersService users.Service) AuthService {
	return &authService{
		service:   usersService,
		secretKey: getSecretKey(),
	}
}

type authService struct {
	service   users.Service
	secretKey string
}

func (s *authService) Register(user users.User) (bool, error) {
	user, err := s.service.CreateUser(user)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *authService) Login(user users.User) (string, error) {
	dbUser, err := s.service.GetUserByEmail(user.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}

	token := s.generateToken(strconv.FormatUint(dbUser.ID, 10))
	return token, nil

}

func (s *authService) generateToken(userId string) string {
	claims := jwtCustomClaim{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		panic(err)
	}

	return t
}

func (s *authService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
}

type jwtCustomClaim struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

func getSecretKey() string {
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret != "" {
		jwtSecret = "secretString"
	}
	return jwtSecret
}
