package security

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTStruct struct {
	Name            string `json:"username"`
	CreatedAt       string `json:"createdAt"`
	PublicUserId    string `json:"public_user_id"`
	PublicSessionId string `json:"public_session_id"`
	SessionSecret   string `json:"session_secret"`
	Email           string `json:"email"`
	InitalizedAt    int64  `json:"iat"`
}

type JwtOptions struct {
	InitalizedAt int64
}

func JWTEncode(jwtStruct JWTStruct, jwtKey string, opts *JwtOptions) (string, error) {
	jwtSecretByte := []byte(jwtKey)

	iat := time.Now().UnixMilli()
	if opts != nil {
		iat = opts.InitalizedAt
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":          jwtStruct.Name,
		"createdAt":         jwtStruct.CreatedAt,
		"public_user_id":    jwtStruct.PublicUserId,
		"public_session_id": jwtStruct.PublicSessionId,
		"session_secret":    jwtStruct.SessionSecret,
		"email":             jwtStruct.Email,
		"iat":               iat,
	})

	tokenString, err := token.SignedString(jwtSecretByte)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func JWTDecode(jwtToken string, jwtKey string) (JWTStruct, error) {
	jwtSecretByte := []byte(jwtKey)

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretByte, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return JWTStruct{}, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		rawJWTToken := token.Claims.(jwt.MapClaims)

		var jwtStruct JWTStruct
		jsonbody, err := json.Marshal(rawJWTToken)
		if err != nil {
			return JWTStruct{}, err
		}

		err = json.Unmarshal(jsonbody, &jwtStruct)
		if err != nil {
			return JWTStruct{}, err
		}

		return jwtStruct, nil
	} else {
		return JWTStruct{}, err
	}
}
