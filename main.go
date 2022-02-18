package main

import (
    "context"
    "fmt"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/bxcodec/faker/v3"
    "strings"
)

// S3CreateBucketAPI defines the interface for the CreateBucket function.
// We use this interface to test the function using a mocked service.
type S3CreateBucketAPI interface {
    CreateBucket(
        ctx context.Context,
        params *s3.CreateBucketInput,
        optFns ...func(*s3.Options),
    ) (*s3.CreateBucketOutput, error)
}

// MakeBucket creates an Amazon Simple Storage Service (Amazon S3) bucket.
// Inputs:
//     c is the context of the method call, which includes the AWS Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a CreateBucketOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to CreateBucket.
func MakeBucket(c context.Context, api S3CreateBucketAPI, input *s3.CreateBucketInput,
) (*s3.CreateBucketOutput, error) {
    return api.CreateBucket(c, input)
}
gc-b
func main() {
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
            func(service, region string, options ...interface{}) (aws.Endpoint, error) {
                return aws.Endpoint{
                    URL: "http://localhost:4566",
                }, nil
            })),
    )
    if err != nil {
        panic("configuration error, " + err.Error())
    }
    client := s3.NewFromConfig(cfg,
        func(options *s3.Options) {
            options.UsePathStyle = true
            options.EndpointOptions.DisableHTTPS = true
        })
    input := &s3.ListBucketsInput{}
    newBucketName := strings.ToLower(faker.FirstName())
    createBucket(*client, newBucketName)
    listBuckets(*client, *input)
}

func listBuckets(client s3.Client, input s3.ListBucketsInput) {
    buckets, err := client.ListBuckets(context.Background(), &input)
    if err != nil {
        fmt.Println("Could not list ")
    }
    for i, b := range buckets.Buckets {
        fmt.Println(i, aws.ToString(b.Name))
    }
}

func createBucket(client s3.Client, bucketName string) {
    cbi := s3.CreateBucketInput{
        Bucket: aws.String(bucketName),
    }
    _, err := client.CreateBucket(context.TODO(), &cbi)
    if err != nil {
        fmt.Println("Could not create bucketName " + bucketName)
    }
}
