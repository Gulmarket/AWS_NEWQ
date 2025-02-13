package main

import (
    "context"
    "fmt"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/dgrijalva/jwt-go"
    "github.com/gin-contrib/cors"
    "github.com/gulmarket/handlers"
    "github.com/gulmarket/repository"
    "github.com/gulmarket/service"

    "github.com/gin-gonic/gin"
    "log"
    "os"
    "strconv"
    "time"
)

func inject(d *dataSources) (*gin.Engine, error) {

    log.Println("Injecting data sources")

    plantationRepository := repository.NewPlantationRepository(d.DB)
    tokenRepository := repository.NewTokenRepository(d.RedisClient)

    proUserRepository := repository.NewProUserRepository(d.DB)

    awsRepository := repository.NewAwsRepository(d.DB)

    plantationService := service.NewPlantationService(&service.PSConfig{
        PlantationRepository: plantationRepository,
        ProUserRepository:    proUserRepository,
    })

    proUserService := service.NewProUserService(&service.PUSConfig{
        ProUserRepository: proUserRepository,
    })

    awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))

    if err != nil {
        log.Printf("Error Loading AWS config")
    }

    awsService := service.NewAWSService(&service.ASConfig{
        S3Client:      s3.NewFromConfig(awsConfig),
        AWSRepository: awsRepository,
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

    router.Use(cors.New(cors.Config{
        AllowAllOrigins:  true,
        AllowCredentials: true,
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},

        //AllowOrigins:  []string{"http://pro.gulmarket.com", "http://localhost", "http://localhost:4251"},
        AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        //AllowHeaders:  []string{"Origin", "Content-Length", "Content-Type"},
        ExposeHeaders: []string{"Content-Length"},
        //AllowOriginFunc: func(origin string) bool {
        //  return origin == "http://pro.gulmarket.com" || origin == "http://localhost" || origin == "http://localhost:3000"
        //},
        MaxAge: 12 * time.Hour,
    }))

    baseURL := os.Getenv("GULMARKET_API_URL")

    handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
    ht, err := strconv.ParseInt(handlerTimeout, 0, 64)
    if err != nil {
        return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
    }

    maxBodyBytes := os.Getenv("MAX_BODY_BYTES")
    mbb, err := strconv.ParseInt(maxBodyBytes, 0, 64)
    if err != nil {
        return nil, fmt.Errorf("could not parse MAX_BODY_BYTES as int: %w", err)
    }

    handlers.NewHandler(&handlers.Config{
        R:                 router,
        PlantationService: plantationService,
        ProUserService:    proUserService,
        TokenService:      tokenService,
        AWSService:        awsService,
        BaseURL:           baseURL,
        TimeoutDuration:   time.Duration(time.Duration(ht) * time.Second),
        MaxBodyBytes:      mbb,
    })

    return router, nil
}
