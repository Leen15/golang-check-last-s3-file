package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/caarlos0/env"
    "fmt"
    "os"
    "time"
    "strings"
)

type config struct {
    Bucket                  string      `env:"BUCKET"`
    Prefix                  string      `env:"PATH_PREFIX"`
    Timezone                string      `env:"TIMEZONE"`
    OlderThanMinutes        int         `env:"OLDER_THAN_MINUTES"`
    SmallerThanMegabytes    int64       `env:"SMALLER_THAN_MEGABYTES"`
    AwsAccessKeyID          string      `env:"AWS_ACCESS_KEY_ID"`
    AwsSecretAccessKey      string      `env:"AWS_SECRET_ACCESS_KEY"`
    AwsRegion               string      `env:"AWS_REGION"`
}

func main() {
    
    cfg := config{}
    if err := env.Parse(&cfg); err != nil {
        exitErrorf("%+v\n", err)
    }

    fmt.Println("Check last file in:")
    fmt.Println("Region:         ", cfg.AwsRegion)
    fmt.Println("Bucket:         ", cfg.Bucket)
    fmt.Println("Prefix:         ", cfg.Prefix)
    fmt.Println("Older than:     ", cfg.OlderThanMinutes, "minutes")
    fmt.Println("Smaller than:   ", cfg.SmallerThanMegabytes, "MB")

    // utc life
    loc, err := time.LoadLocation(cfg.Timezone)
    if err != nil {
        panic(err)
    }

    sess, err := session.NewSession(&aws.Config{
        Region:      aws.String(cfg.AwsRegion),
        Credentials: credentials.NewStaticCredentials(cfg.AwsAccessKeyID, cfg.AwsSecretAccessKey, ""),
    })

    // Create S3 service client
    svc := s3.New(sess)
    
    // Get the list of items
    resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(cfg.Bucket), Prefix: aws.String(cfg.Prefix)})
    if err != nil {
        exitErrorf("[ERROR] Unable to list items in bucket %q, %v", cfg.Bucket, err)
    }

    if (len(resp.Contents) == 0) {
        fmt.Println("[ERROR] No files in bucket. ")
        os.Exit(0)
    }

    mostRecentObj := *resp.Contents[0]
    for _, item := range resp.Contents {
        if (item.LastModified.After(*mostRecentObj.LastModified)) {
            mostRecentObj = *item
        } 
    }

    fmt.Println("Files in bucket:", len(resp.Contents))
    fmt.Println("")

    keyArray := strings.Split(*mostRecentObj.Key, "/")
    mostRecentObjName := keyArray[len(keyArray)-1]
    mostRecentObjDate := (*mostRecentObj.LastModified).In(loc)
    mostRecentObjSize := *mostRecentObj.Size  / 1024 / 1024

    fmt.Println("Most recent file is: ")
    fmt.Println("Name:         ", mostRecentObjName)
    fmt.Println("modified at:  ", mostRecentObjDate)
    fmt.Println("Size:         ", mostRecentObjSize , "MB")
    fmt.Println("")

    minOlderTime := time.Now().In(loc).Add(time.Minute * time.Duration(-1 * cfg.OlderThanMinutes))
    diffTime := mostRecentObjDate.Sub(minOlderTime);

    error := false
    if (diffTime.Minutes() < 0) {
        fmt.Println("[ERROR] The file is older than max allowed. (", diffTime * -1, "ago )" )
        error = true
    }

    if (cfg.SmallerThanMegabytes > 0 && *mostRecentObj.Size  / 1024 / 1024  < cfg.SmallerThanMegabytes) {
        fmt.Println("[ERROR] The file is smaller than min allowed. (", mostRecentObjSize, "MB vs", cfg.SmallerThanMegabytes, "MB )" )
        error = true
    }

    if (error) {
        os.Exit(0)
    }

    fmt.Println("[SUCCESS] The file is OK." )	
}

func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}