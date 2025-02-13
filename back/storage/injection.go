package storage

import (
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "github.com/gulmarket/handlers"
    "github.com/gulmarket/repository"
    "github.com/gulmarket/service"
    "log"
    "os"
    "strconv"
    "time"
)

func Inject(d *DataSources) (*gin.Engine, error) {

    log.Println("Injecting data sources")

    plantationRepository := repository.NewPlantationRepository(d.DB)
    tokenRepository := repository.NewTokenRepository(d.RedisClient)

    plantationService := service.NewPlantationService(&service.PSConfig{
        PlantationRepository: plantationRepository,
    })

    proUserRepository := repository.NewProUserRepository(d.DB)

    proUserService := service.NewProUserService(&service.PUSConfig{
        ProUserRepository: proUserRepository,
    })

    // load rsa keys
    privKeyFile := os.Getenv("PRIV_KEY_FILE")
    priv, err := os.ReadFile(privKeyFile)

    if err != nil {
        return nil, fmt.Errorf("could not read private key pem file: %w", err)
    }

    privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)

    if err != nil {
        return nil, fmt.Errorf("could not parse private key: %w", err)
    }

    pubKeyFile := os.Getenv("PUB_KEY_FILE")
    pub, err := os.ReadFile(pubKeyFile)

    if err != nil {
        return nil, fmt.Errorf("could not read public key pem file: %w", err)
    }

    pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)

    if err != nil {
        return nil, fmt.Errorf("could not parse public key: %w", err)
    }

    // load refresh token secret from env variable
    refreshSecret := os.Getenv("REFRESH_SECRET")

    // load expiration lengths from env variables and parse as int
    idTokenExp := os.Getenv("ID_TOKEN_EXP")
    refreshTokenExp := os.Getenv("REFRESH_TOKEN_EXP")

    idExp, err := strconv.ParseInt(idTokenExp, 0, 64)
    if err != nil {
        return nil, fmt.Errorf("could not parse ID_TOKEN_EXP as int: %w", err)
    }

    refreshExp, err := strconv.ParseInt(refreshTokenExp, 0, 64)
    if err != nil {
        return nil, fmt.Errorf("could not parse REFRESH_TOKEN_EXP as int: %w", err)
    }

    tokenService := service.NewTokenService(&service.TSConfig{
        TokenRepository:       tokenRepository,
        PrivKey:               privKey,
        PubKey:                pubKey,
        RefreshSecret:         refreshSecret,
        IDExpirationSecs:      idExp,
        RefreshExpirationSecs: refreshExp,
    })

    router := gin.Default()
    baseURL := os.Getenv("GULMARKET_API_URL")

    handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
    ht, err := strconv.ParseInt(handlerTimeout, 0, 64)
    if err != nil {
        return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
    }

    handlers.NewHandler(&handlers.Config{
        R:                 router,
        PlantationService: plantationService,
        ProUserService:    proUserService,
        TokenService:      tokenService,
        BaseURL:           baseURL,
        TimeoutDuration:   time.Duration(time.Duration(ht) * time.Second),
    })

    return router, nil
}
