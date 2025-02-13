package service

import (
    "context"
    "crypto/rsa"
    "github.com/gulmarket/model"
    "github.com/gulmarket/model/apperrors"

    "log"
    "strconv"
)

type tokenService struct {
    TokenRepository       model.TokenRepository
    PrivKey               *rsa.PrivateKey
    PubKey                *rsa.PublicKey
    RefreshSecret         string
    IDExpirationSecs      int64
    RefreshExpirationSecs int64
}

type TSConfig struct {
    TokenRepository       model.TokenRepository
    PrivKey               *rsa.PrivateKey
    PubKey                *rsa.PublicKey
    RefreshSecret         string
    IDExpirationSecs      int64
    RefreshExpirationSecs int64
}

func NewTokenService(c *TSConfig) model.TokenService {
    return &tokenService{
        TokenRepository:       c.TokenRepository,
        PrivKey:               c.PrivKey,
        PubKey:                c.PubKey,
        RefreshSecret:         c.RefreshSecret,
        IDExpirationSecs:      c.IDExpirationSecs,
        RefreshExpirationSecs: c.RefreshExpirationSecs,
    }
}

func (s *tokenService) NewPairFromPlantation(ctx context.Context, p *model.Plantation, prevTokenID string) (*model.TokenPair, error) {
    id := strconv.Itoa(p.Id)
    if prevTokenID != "" {
        if err := s.TokenRepository.DeleteRefreshToken(ctx, id, prevTokenID); err != nil {
            log.Printf("Could not delete previous refreshToken for id: %v, tokenID: %v\n", id, prevTokenID)

            return nil, err
        }
    }

    idToken, err := generateIDToken(p, s.PrivKey, s.IDExpirationSecs)

    if err != nil {
        log.Printf("Error generating idToken for id: %v. Error: %v\n", id, err.Error())
        return nil, apperrors.NewInternal()
    }

    refreshToken, err := generateRefreshToken(id, s.RefreshSecret, s.RefreshExpirationSecs)

    if err != nil {
        log.Printf("Error generating refreshToken for id: %v. Error: %v\n", id, err.Error())
        return nil, apperrors.NewInternal()
    }

    if err := s.TokenRepository.SetRefreshToken(ctx, id, refreshToken.ID, refreshToken.ExpiresIn); err != nil {
        log.Printf("Error storing tokenID for id: %v. Error: %v\n", id, err.Error())
        return nil, apperrors.NewInternal()
    }

    return &model.TokenPair{
        IDToken:      model.IDToken{SS: idToken},
        RefreshToken: model.RefreshToken{SS: refreshToken.SS, ID: refreshToken.ID},
    }, nil
}

func (s *tokenService) NewPairFromProUser(ctx context.Context, p *model.ProUser, prevTokenID string) (*model.TokenPair, error) {
    id := strconv.Itoa(p.Id)
    if prevTokenID != "" {
        if err := s.TokenRepository.DeleteRefreshToken(ctx, id, prevTokenID); err != nil {
            log.Printf("Could not delete previous refreshToken for id: %v, tokenID: %v\n", id, prevTokenID)

            return nil, err
        }
    }

    idToken, err := generateProUserIdToken(p, s.PrivKey, s.IDExpirationSecs)

    if err != nil {
        log.Printf("Error generating idToken for id: %v. Error: %v\n", id, err.Error())
        return nil, apperrors.NewInternal()
    }

    refreshToken, err := generateRefreshToken(id, s.RefreshSecret, s.RefreshExpirationSecs)

    if err != nil {
        log.Printf("Error generating refreshToken for id: %v. Error: %v\n", id, err.Error())
        return nil, apperrors.NewInternal()
    }

    if err := s.TokenRepository.SetRefreshToken(ctx, id, refreshToken.ID, refreshToken.ExpiresIn); err != nil {
        log.Printf("Error storing tokenID for id: %v. Error: %v\n", id, err.Error())
        return nil, apperrors.NewInternal()
    }

    return &model.TokenPair{
        IDToken:      model.IDToken{SS: idToken},
        RefreshToken: model.RefreshToken{SS: refreshToken.SS, ID: refreshToken.ID},
    }, nil
}

func (s *tokenService) Signout(ctx context.Context, id int) error {
    uid := strconv.Itoa(id)
    return s.TokenRepository.DeletePlantationRefreshTokens(ctx, uid)
}

func (s *tokenService) ValidateIDToken(tokenString string) (*model.Plantation, error) {
    claims, err := validateIDToken(tokenString, s.PubKey) // uses public RSA key

    if err != nil {
        log.Printf("Unable to validate or parse idToken - Error: %v\n", err)
        return nil, apperrors.NewAuthorization("Unable to verify plantation from idToken")
    }

    return claims.Plantation, nil
}

func (s *tokenService) ValidateProUserIDToken(tokenString string) (*model.ProUser, error) {
    claims, err := validateProUserIdToken(tokenString, s.PubKey)

    if err != nil {
        log.Printf("Unable to validate or parse idToken - Error: %v\n", err)
        return nil, apperrors.NewAuthorization("Unable to verify user from idToken")
    }

    return claims.ProUser, nil
}

func (s *tokenService) ValidateRefreshToken(tokenString string) (*model.RefreshToken, error) {
    claims, err := validateRefreshToken(tokenString, s.RefreshSecret)

    if err != nil {
        log.Printf("Unable to validate or parse refreshToken for token string: %s\n%v\n", tokenString, err)
        return nil, apperrors.NewAuthorization("Unable to verify plantation from refresh token")
    }

    if err != nil {
        log.Printf("Claims ID could not be parsed: %s\n%v\n", claims.Id, err)
        return nil, apperrors.NewAuthorization("Unable to verify plantation from refresh token")
    }

    return &model.RefreshToken{
        SS: tokenString,
        ID: claims.ID,
    }, nil
}
