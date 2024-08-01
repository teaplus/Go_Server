package ultils

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	User string `json:"user_id"`
	jwt.StandardClaims
}

func ValidateToken(tokenString, secretKey string) (*Claims, error) {
	fmt.Println("Key:", secretKey)
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		return nil, err
	}
	
	fmt.Printf("Token claims: %v\n", token.Claims)

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	if claimsMap, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Claims map:", claimsMap)
		userID, ok := claimsMap["user_id"].(string)
		if !ok {
			return nil, fmt.Errorf("unable to cast user_id to string")
		}
		exp, ok := claimsMap["exp"].(float64)
		if !ok {
			return nil, fmt.Errorf("unable to cast exp to float64")
		}
		return &Claims{
			User: userID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: int64(exp),
			},
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func CreateTokenpair(payload Claims, publicKey, privateKey string) (map[string]string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	accessTokenString, err := accessToken.SignedString([]byte(publicKey))
	if err != nil {
		return nil, err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	refreshTokenString, err := refreshToken.SignedString([]byte(privateKey))
	if err != nil {
		return nil, err
	}
	access, err := jwt.Parse(accessTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(publicKey), nil
	})

	if err != nil {
		return nil, err
	}

	fmt.Println("accc", access)
	fmt.Println("publickey", publicKey)

	return map[string]string{
		"accessToken":  accessTokenString,
		"refreshToken": refreshTokenString,
	}, nil
}
