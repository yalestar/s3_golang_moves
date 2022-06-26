✅ get docker-compose to run S3 
✅ create and list a bucket via CLI?  
✅ create and list bucket via Go  
✅ add a file to localstack S3   
✅ get that file from localstack S3 

Run the create and list with that interfaced version  
like shown here https://aws.github.io/aws-sdk-go-v2/docs/unit-testing/

===

#### Sample commands:
aws --endpoint-url=http://localhost:4566 s3api list-buckets
aws --endpoint-url=http://localhost:4566 s3 ls s3://hoobyak
