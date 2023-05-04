package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const bucketName = "tokatu4561-test-bucket-1234"
const region =  "ap-north-east1"

func main() {
	ctx := context.Background()
	s3Client, err := initS3Client(ctx)
	if err != nil {
		fmt.Printf("inits3Clitnent error: %s", err)
	}

	if err = createS3BucketIfNotExist(ctx, s3Client); err != nil {
		fmt.Printf("create s3 bucket error: %s", err)
		os.Exit(1)
	}

	if err = uploadToS3Bucket(ctx, s3Client); err != nil {
		fmt.Printf("update s3 bucket error: %s", err)
		os.Exit(1)
	}

	fmt.Printf("upload complete")
}

func initS3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }
	
	return s3.NewFromConfig(cfg), nil
}

// s3バケットが存在していなければ作成
func createS3BucketIfNotExist(ctx context.Context, client *s3.Client) error {
	allBuckets, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("ListBuckets error: %s", err)
	}

	found := false
	for _, b := range allBuckets.Buckets {
		if *b.Name == bucketName {
			found = true
		}
	}

	if !found {
		_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: region,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func uploadToS3Bucket(ctx context.Context, client *s3.Client) error {
	file, err := ioutil.ReadFile("test.txt")
	if err != nil {
		return err
	}

	uploader := manager.NewUploader(client)

	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Key: aws.String("text.txt"),
		Bucket: aws.String(bucketName),
		Body: bytes.NewReader(file),
	})
	if err != nil {
		return err
	}

	return nil
}