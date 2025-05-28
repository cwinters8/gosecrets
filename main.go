package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(os.Getenv("AWS_PROFILE")))
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	svc := secretsmanager.NewFromConfig(cfg)

	secrets := map[string]map[string]string{}
	secretsFile, err := os.ReadFile(os.Getenv("SECRETS_FILE_PATH"))
	if err != nil {
		log.Fatalf("failed to read secrets file: %v", err)
	}

	if err := json.Unmarshal(secretsFile, &secrets); err != nil {
		log.Fatalf("failed to unmarshal secrets: %v", err)
	}

	for k, v := range secrets {
		secretString, err := json.Marshal(v)
		if err != nil {
			log.Fatalf("failed to marshal secret: %v", err)
		}
		if _, err := svc.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(k),
			SecretString: aws.String(string(secretString)),
		}); err != nil {
			log.Fatalf("failed to create secret: %v", err)
		}
	}
}
