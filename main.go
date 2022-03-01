package main

import (
    "bufio"
    "context"
    "errors"
    "fmt"
    v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
    awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
    "io"
    "log"
    "net/http"
    "os"
    "time"
    
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3CreateBucketAPI defines the interface for the CreateBucket function.
// We use this interface to test the function using a mocked service.
type S3CreateBucketAPI interface {
    CreateBucket(
        ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options),
    ) (*s3.CreateBucketOutput, error)
}

func MakeBucket(
    c context.Context, api S3CreateBucketAPI, input *s3.CreateBucketInput,
) (*s3.CreateBucketOutput, error) {
    return api.CreateBucket(c, input)
}

func main() {
    cfg, err := config.LoadDefaultConfig(
        context.TODO(),
        config.WithEndpointResolverWithOptions(
            aws.EndpointResolverWithOptionsFunc(
                func(service, region string, options ...interface{}) (
                    aws.Endpoint, error,
                ) {
                    return aws.Endpoint{
                        URL: "http://localhost:4566",
                    }, nil
                },
            ),
        ),
    )
    if err != nil {
        log.Fatal("configuration error, " + err.Error())
    }
    
    client := s3.NewFromConfig(
        cfg,
        func(options *s3.Options) {
            options.UsePathStyle = true
            options.EndpointOptions.DisableHTTPS = true
        },
    )
    input := &s3.ListBucketsInput{}
    // newBukit := strings.ToLower(faker.FirstName())
    newBukit := "hoobyak"
    exists, err := bucketExists(*client, newBukit)
    if !exists {
        createBucket(*client, newBukit)
    }
    listBuckets(*client, *input)
    thing, err := putObject(*client, newBukit)
    if err != nil {
        log.Println(err)
    }
    fmt.Println(thing)
    time.Sleep(time.Second * 4)
    
    fud, err := getObject(*client, "hoobyak", "doogie.txt")
    if err != nil {
        log.Println(err)
    }
    all, err := io.ReadAll(fud.Body)
    
    if all == nil {
        log.Fatal("all is nothing. 0+2=1")
    }
    if err != nil {
        return
    }
    fmt.Println(string(all))
}

func getObject(client s3.Client, bucket, key string) (*s3.GetObjectOutput, error) {
    goi := s3.GetObjectInput{Bucket: &bucket, Key: &key}
    obj, err := client.GetObject(context.TODO(), &goi)
    
    if err != nil {
        log.Println(err)
        return nil, err
    }
    
    return obj, nil
}

func putObject(client s3.Client, bucket string) (string, error) {
    dookFile := "/Users/r622233/dev/EXAMPLES/s3mock/dook.txt"
    file, _ := os.Open(dookFile)
    f := bufio.NewReader(file)
    input := s3.PutObjectInput{
        Bucket: aws.String(bucket),
        Body:   f,
        Key:    aws.String("doogie.txt"),
    }
    
    _, err := client.PutObject(
        context.Background(), &input, s3.WithAPIOptions(
            v4.SwapComputePayloadSHA256ForUnsignedPayloadMiddleware,
        ),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    return "ok", nil
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
        log.Fatal(err)
    }
}

func bucketExists(client s3.Client, bucketName string) (bool, error) {
    _, err := client.HeadBucket(
        context.TODO(), &s3.HeadBucketInput{
            Bucket: aws.String(bucketName),
        },
    )
    if err != nil {
        var respError *awshttp.ResponseError
        if errors.As(
            err, &respError,
        ) && respError.ResponseError.HTTPStatusCode() == http.StatusNotFound {
            return false, nil
        }
        return false, err
    }
    return true, nil
}
