package service

import (
    "crypto/rsa"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/gulmarket/model"
    "log"
    "time"
)

type idTokenCustomClaims struct {
    Plantation *model.Plantation `json:"plantation"`
    jwt.StandardClaims
}

type proUserIdTokenCustomClaims struct {
    ProUser *model.ProUser `json:"pro_user"`
    jwt.StandardClaims
}

func generateIDToken(p *model.Plantation, key *rsa.PrivateKey, exp int64) (string, error) {
    unixTime := time.Now().Unix()
    tokenExp := unixTime + exp

    claims := idTokenCustomClaims{
        Plantation: p,
        StandardClaims: jwt.StandardClaims{
            IssuedAt:  unixTime,
            ExpiresAt: tokenExp,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    ss, err := token.SignedString(key)

    if err != nil {
        log.Println("Failed to sign id token string")
        return "", err
    }

    return ss, nil
}

func generateProUserIdToken(p *model.ProUser, key *rsa.PrivateKey, exp int64) (string, error) {
    unixTime := time.Now().Unix()
    tokenExp := unixTime + exp

    claims := proUserIdTokenCustomClaims{
        ProUser: p,
        StandardClaims: jwt.StandardClaims{
            IssuedAt:  unixTime,
            ExpiresAt: tokenExp,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    ss, err := token.SignedString(key)

    if err != nil {
        log.Println("Failed to sign id token string")
        return "", err
    }

    return ss, nil
}

type refreshTokenData struct {
    SS        string
    ID        string
    ExpiresIn time.Duration
}

type refreshTokenCustomClaims struct {
    ID string `json:"id"`
    jwt.StandardClaims
}

func generateRefreshToken(id string, key string, exp int64) (*refreshTokenData, error) {
    currentTime := time.Now()
    tokenExp := currentTime.Add(time.Duration(exp) * time.Second)

    claims := refreshTokenCustomClaims{
        ID: id,
        StandardClaims: jwt.StandardClaims{
            IssuedAt:  currentTime.Unix(),
            ExpiresAt: tokenExp.Unix(),
            Id:        id,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    ss, err := token.SignedString([]byte(key))

    if err != nil {
        log.Println("Failed to sign refresh token string")
        return nil, err
    }

    return &refreshTokenData{
        SS:        ss,
        ID:        id,
        ExpiresIn: tokenExp.Sub(currentTime),
    }, nil
}

func validateIDToken(tokenString string, key *rsa.PublicKey) (*idTokenCustomClaims, error) {
    claims := &idTokenCustomClaims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return key, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("ID token is invalid")
    }

    claims, ok := token.Claims.(*idTokenCustomClaims)

    if !ok {
        return nil, fmt.Errorf("ID token valid but couldn't parse claims")
    }

    return claims, nil
}

func validateProUserIdToken(tokenString string, key *rsa.PublicKey) (*proUserIdTokenCustomClaims, error) {
    claims := &proUserIdTokenCustomClaims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return key, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("ID token is invalid")
    }

    claims, ok := token.Claims.(*proUserIdTokenCustomClaims)

    if !ok {
        return nil, fmt.Errorf("ID token valid but couldn't parse claims")
    }

    return claims, nil
}

func validateRefreshToken(tokenString string, key string) (*refreshTokenCustomClaims, error) {
    claims := &refreshTokenCustomClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(key), nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("Refresh token is invalid")
    }

    claims, ok := token.Claims.(*refreshTokenCustomClaims)

    if !ok {
        return nil, fmt.Errorf("Refresh token valid but couldn't parse claims")
    }

    return claims, nil
}
