package main

import (
    "bufio"
    "context"
    "errors"
    "fmt"
    v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
    awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
    "github.com/bxcodec/faker/v3"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
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
    newBucketName := faker.FirstName()
    
    cfg, err := getS3Config()
    if err != nil {
        log.Fatal("NO CONNECTO" + err.Error())
    }
    
    client := getS3Client(cfg)
    input := &s3.ListBucketsInput{}
    
    exists, err := bucketExists(*client, newBucketName)
    if !exists {
        createBucket(*client, newBucketName)
    }
    listBuckets(*client, *input)
    // later TODO: make this wait for a signal or a channel response
    time.Sleep(time.Second * 2)
    
    fileToUpload, err := createDookFile()
    _, _ = putObject(*client, newBucketName, fileToUpload)
    
    fud, err := getObject(*client, newBucketName, "doogie.txt")
    if err != nil {
        log.Println(err)
    }
    all, err := io.ReadAll(fud.Body)
    
    if all == nil {
        log.Fatal("All is nothing. Freedom is slavery. 0+2=1")
    }
    fmt.Println(string(all))
}

func getS3Client(cfg aws.Config) *s3.Client {
    client := s3.NewFromConfig(
        cfg,
        func(options *s3.Options) {
            options.UsePathStyle = true
            options.EndpointOptions.DisableHTTPS = true
        },
    )
    return client
}

func getS3Config() (aws.Config, error) {
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
    return cfg, err
}

type S3GetObjAPI interface {
    GetObject(
        ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options),
    ) (*s3.GetObjectOutput, error)
}

func GetObjectTheRadWay(api S3GetObjAPI, bucket, key string) ([]byte, error) {
    goi := s3.GetObjectInput{Bucket: &bucket, Key: &key}
    obj, err := api.GetObject(context.TODO(), &goi)
    if err != nil {
        log.Println(err)
    }
    defer obj.Body.Close()
    
    return io.ReadAll(obj.Body)
}

//goland:noinspection ALL
func createDookFile() (string, error) {
    log.Println("------- Fixin to create text file ---------")
    pwd, _ := os.Getwd()
    fp := filepath.Join(pwd, "dook.txt")
    file, err := os.Create(fp)
    if err != nil {
        log.Println(err)
        return "", nil
    }
    
    defer file.Close()
    
    for i := 0; i < 10; i++ {
        file.WriteString(faker.Name() + "\n")
    }
    
    file.Sync()
    
    return file.Name(), nil
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

func putObject(client s3.Client, bucket, fileName string) (string, error) {
    file, _ := os.Open(fileName)
    f := bufio.NewReader(file)
    bfn := filepath.Base(fileName)
    log.Printf("------- Fixin to put %s to %s ---------\n", bfn, bucket)
    
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
        log.Println("Could not list ")
    }
    
    log.Println("------- Listing Buckets 4 U --------")
    for i, b := range buckets.Buckets {
        log.Println(i, aws.ToString(b.Name))
    }
}

func createBucket(client s3.Client, bucketName string) {
    cbi := s3.CreateBucketInput{
        Bucket: aws.String(bucketName),
    }
    log.Printf("------- Fixin to create bucket %s --------\n", bucketName)
    _, err := client.CreateBucket(context.TODO(), &cbi)
    if err != nil {
        log.Println("Could not create bucketName " + bucketName)
        log.Fatal(err)
    }
}

func bucketExists(client s3.Client, bucketName string) (bool, error) {
    _, err := client.HeadBucket(
        context.TODO(), &s3.HeadBucketInput{
            Bucket: aws.String(bucketName),
        },
    )
    
    log.Printf("------- Checking if %s exists -------\n", bucketName)
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
