package repository

import (
    "context"
    "fmt"
    "github.com/gulmarket/model"
    "github.com/gulmarket/model/apperrors"
    "github.com/redis/go-redis/v9"
    "log"
    "time"
)

type redisTokenRepository struct {
    Redis *redis.Client
}

func NewTokenRepository(redisClient *redis.Client) model.TokenRepository {
    return &redisTokenRepository{
        Redis: redisClient,
    }
}

func (r *redisTokenRepository) SetRefreshToken(ctx context.Context, id string, tokenID string, expiresIn time.Duration) error {

    key := fmt.Sprintf("%s:%s", id, tokenID)
    if err := r.Redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
        log.Printf("Could not SET refresh token to redis for id/tokenID: %s/%s: %v\n", id, tokenID, err)
        return apperrors.NewInternal()
    }
    return nil
}

func (r *redisTokenRepository) DeleteRefreshToken(ctx context.Context, ID string, tokenID string) error {
    key := fmt.Sprintf("%s:%s", ID, tokenID)

    result := r.Redis.Del(ctx, key)

    if err := result.Err(); err != nil {
        log.Printf("Could not delete refresh token to redis for id/tokenID: %s/%s: %v\n", ID, tokenID, err)
        return apperrors.NewInternal()
    }

    if result.Val() < 1 {
        log.Printf("Refresh token to redis for ID/tokenID: %s/%s does not exist\n", ID, tokenID)
        return apperrors.NewAuthorization("Invalid refresh token")
    }

    return nil
}

func (r *redisTokenRepository) DeletePlantationRefreshTokens(ctx context.Context, id string) error {
    pattern := fmt.Sprintf("%s*", id)

    iter := r.Redis.Scan(ctx, 0, pattern, 5).Iterator()
    failCount := 0

    for iter.Next(ctx) {
        if err := r.Redis.Del(ctx, iter.Val()).Err(); err != nil {
            log.Printf("Failed to delete refresh token: %s\n", iter.Val())
            failCount++
        }
    }

    // check last value
    if err := iter.Err(); err != nil {
        log.Printf("Failed to delete refresh token: %s\n", iter.Val())
    }

    if failCount > 0 {
        return apperrors.NewInternal()
    }

    return nil
}
