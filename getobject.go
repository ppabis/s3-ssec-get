package main

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

/*
	Recursively download all the objects in a bucket and prefix with SSE-C decryption.
*/
func RecursiveGetObject(cfg aws.Config, bucketName string, prefix string, ssecKey string, output string) {
	client := s3.NewFromConfig(cfg)

	// Get all the objects in the bucket
	objects, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucketName,
		Prefix: &prefix,
	})

	if err != nil {
		log.Fatalf("[ERROR] listing objects in bucket %s: %s\n", bucketName, err)
	}

	// Iterate over all the objects
	for _, object := range objects.Contents {
		log.Default().Printf("getting object %s", *object.Key)
		GetObject(client, bucketName, prefix, *object.Key, ssecKey, output)
	}
}

/*
	Specify paths for the object to download.

	So it is given the key of the object in the bucket. From this key is
	substracted the prefix. Then this result is appended to the output directory
	path. All the necessary directories are created before copying the object.
*/
func GetObject(client *s3.Client, bucketName string, prefix string, key string, ssecKey string, output string) {
	// Remove the prefix from the key
	newKey := key[len(prefix):]
	// Append the new key to the output directory
	output = filepath.Join(output, newKey)
	// Check if all the parent directories exist or create them
	dir := filepath.Dir(output)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	err := transferObject(client, bucketName, key, ssecKey, output)
	if err != nil {
		log.Fatalf("[ERROR] transferring object %s: %s\n", key, err)
	}
}

/*
	Download a single object with SSE-C decryption.
*/
func transferObject(client *s3.Client, bucketName string, key string, ssecKey string, output string) error {
	keyDigest := keyMd5(ssecKey)

	object, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket:               &bucketName,
		Key:                  &key,
		SSECustomerAlgorithm: aws.String("AES256"),
		SSECustomerKey:       &ssecKey,
		SSECustomerKeyMD5:    &keyDigest,
	})

	if err != nil {
		return fmt.Errorf("[ERROR] getting object %s: %s", key, err)
	}

	of, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("[ERROR] opening file %s: %s", output, err)
	}
	defer of.Close()

	buf := make([]byte, 1024)
	for {
		n, err := object.Body.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("[ERROR] reading object %s: %s", key, err)
		}

		_, werr := of.Write(buf[:n])
		if werr != nil {
			return fmt.Errorf("[ERROR] writing object %s: %s", key, werr)
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}

/*
 Calculate the MD5 hash of the raw SSE-C key
*/
func keyMd5(ssecKey string) string {
	rawKey, err := base64.StdEncoding.DecodeString(ssecKey)
	if err != nil {
		log.Fatalf("[ERROR] decoding ssecKey: %s\n", err)
		return ""
	}
	hasher := md5.New()
	hasher.Write(rawKey)
	keyHashB64 := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	return keyHashB64
}
