package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		fmt.Printf("[ERROR] unable to load SDK config, %v\n", err)
		os.Exit(2)
	}

	if len(os.Args) < 5 {
		fmt.Println("Usage: s3-ssec-get <bucketName> <prefix> <key> <outputDir>")
		os.Exit(1)
	}

	bucketName := os.Args[1]
	prefix := os.Args[2]
	key := os.Args[3]
	outputDir := os.Args[4]

	RecursiveGetObject(cfg, bucketName, prefix, key, outputDir)
}
