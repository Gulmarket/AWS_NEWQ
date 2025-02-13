package service

import (
    "context"
    "fmt"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/google/uuid"
    "github.com/gulmarket/model"
    "mime/multipart"
)

type awsService struct {
    S3Client      *s3.Client
    AWSRepository model.AWSRepository
}

type ASConfig struct {
    S3Client      *s3.Client
    AWSRepository model.AWSRepository
}

func NewAWSService(c *ASConfig) model.AWSService {
    return &awsService{
        S3Client:      c.S3Client,
        AWSRepository: c.AWSRepository,
    }
}

func (u *awsService) UploadImage(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
    file, err := fileHeader.Open()
    if err != nil {
        return "", fmt.Errorf("failed to open file header: %w", err)
    }
    defer file.Close()

    uniqueID := uuid.New().String()
    filename := fmt.Sprintf("%s-%s", uniqueID, fileHeader.Filename)

    _, err = u.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
        Bucket:      aws.String("gulmarketlogos"),
        Key:         aws.String(filename),
        Body:        file,
        ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
    })
    if err != nil {
        return "", fmt.Errorf("failed to put object in S3: %w", err)
    }

    url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", "gulmarketlogos", filename)
    return url, nil
}
