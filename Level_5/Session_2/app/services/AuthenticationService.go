package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"idstar.com/session9/app/dtos"
	"idstar.com/session9/app/repositories"
	"idstar.com/session9/app/tools"
)

const key = "abcdefghij1234567890"

type AuthenticationService struct {
	userRepository repositories.UserRepository
}

func NewAuthenticationService(userRepository repositories.UserRepository) *AuthenticationService {
	return &AuthenticationService{
		userRepository: userRepository,
	}
}

func (u *AuthenticationService) Login(dto dtos.LoginRequest) (int, *dtos.Token, error) {
	userWithRoles, err := u.userRepository.FindByUsernameOrEmailWithRole(dto.Username)

	if err != nil {
		return 500, nil, err
	}

	if len(*userWithRoles) == 0 {
		return 404, nil, errors.New("user not found")
	}

	aes128 := tools.Aes128{}
	decryptedPassword, err := aes128.Decrypt((*userWithRoles)[0].Password)
	if err != nil {
		return 500, nil, err
	}

	if dto.Password != *decryptedPassword {
		return 401, nil, errors.New("invalid password")
	}

	permissions := ConvertToPermissionsList(*userWithRoles)

	token, err := GenerateToken((*userWithRoles)[0].Username, (*userWithRoles)[0].RoleName, permissions)
	if err != nil {
		return 500, nil, err
	}

	return 200, token, nil
}

func GenerateToken(username string, roleName string, permissions []dtos.Permission) (*dtos.Token, error) {
	iat := jwt.NewNumericDate(time.Now())
	exp := jwt.NewNumericDate(time.Now().Add(time.Millisecond * 3600000))

	claims := &dtos.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "issuer",
			Audience:  []string{"audience01", "audience02"},
			IssuedAt:  iat,
			ExpiresAt: exp,
		},
		Role: dtos.RoleSwagger{
			Name:        roleName,
			Permissions: permissions,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return nil, err
	}

	return &dtos.Token{
		Value:     tokenString,
		IssuedOn:  time.Unix(iat.Unix(), 0),
		ExpiresOn: time.Unix(exp.Unix(), 0),
	}, nil
}

func ConvertToPermissionsList(data []dtos.UserWithClaimsResponse) []dtos.Permission {
	result := []dtos.Permission{}
	for _, r := range data {
		result = append(result, dtos.Permission{
			Module:    r.ModuleName,
			CanCreate: r.CanCreate,
			CanRead:   r.CanRead,
			CanUpdate: r.CanUpdate,
			CanDelete: r.CanDelete,
		})
	}
	return result
}
