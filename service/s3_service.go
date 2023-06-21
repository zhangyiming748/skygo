package service

import (
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/
/**
 * s3client              一个简单的AWS S3的客户端。实现本地文件或是其他对象的上传、下载和删除
 * @package              s3client
 */
type S3Client struct {
	S3Bucket    string // Bucket名称，从hulk可以获取
	S3AccessKey string // AccessKey名称，从hulk可以获取
	S3SecretKey string // SecretKey名称，从hulk可以获取
	S3Region    string
	S3EndPoint  string // S3的域名，从hulk可以获取
}

func NewS3Client() *S3Client {
	s3Config := LoadS3Config()
	return &S3Client{
		S3Bucket:    s3Config.Bucket,
		S3AccessKey: s3Config.AccessKey,
		S3SecretKey: s3Config.SecretKey,
		S3Region:    "us-west-2",
		S3EndPoint:  s3Config.EndPoint,
	}
}

/**
*    上传本地文件到S3
*    params :
*        key S3上保存文件的路径（名字）
*        filename 本地文件名字
 */
func (sc *S3Client) UploadFile(key, filename string) (string, error) {
	creds := credentials.NewStaticCredentials(sc.S3AccessKey, sc.S3SecretKey, "")
	config := &aws.Config{
		Region:           aws.String(sc.S3Region),
		Endpoint:         aws.String(sc.S3EndPoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      creds,
	}
	sess := session.Must(session.NewSession(config))

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(filename)
	if err != nil {
		log.Printf("failed to open file %q, %v", filename, err)
		return "", err
	}

	// Upload the file to S3.
	if _, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(sc.S3Bucket),
		Key:    aws.String(key),
		Body:   f,
	}); err != nil {
		log.Printf("failed to upload file, %v", err)
		return "", err
	}
	path, err := sc.GetSignedUrl(key)
	return path, err
}

/**
*    从S3下载文件到本地并保存
*    params :
*        key S3上需要下载的文件的路径（名字）
*        filename 本地文件名字
 */
func (sc *S3Client) DownloadFile(key, filename string) error {
	creds := credentials.NewStaticCredentials(sc.S3AccessKey, sc.S3SecretKey, "")
	config := &aws.Config{
		Region:           aws.String(sc.S3Region),
		Endpoint:         aws.String(sc.S3EndPoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      creds,
	}
	sess := session.Must(session.NewSession(config))

	// Create an uploader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Create a file to write the S3 Object contents to.
	f, err := os.Create(filename)
	if err != nil {
		log.Printf("failed to create file %q, %v", filename, err)
		return err
	}

	// Write the contents of S3 Object to the file
	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(sc.S3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Printf("failed to download file, %v", err)
		return err
	}

	return nil
}

func (sc *S3Client) GetSignedUrl(key string) (string, error) {
	creds := credentials.NewStaticCredentials(sc.S3AccessKey, sc.S3SecretKey, "")
	config := &aws.Config{
		Region:           aws.String(sc.S3Region),
		Endpoint:         aws.String(sc.S3EndPoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      creds,
	}
	sess := session.Must(session.NewSession(config))

	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(sc.S3Bucket),
		Key:    aws.String(key),
	})
	req.Time = time.Now()
	// s3有效期最大7天
	urlStr, err := req.Presign(168 * time.Hour)

	return urlStr, err
}
