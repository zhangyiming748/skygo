package orm_mongo

import (
	"bytes"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 文件上传
func GridFSUpload(collectionName string, filename string, fileContent []byte) (fileId string, err error) {
	bucket, err := gridfs.NewBucket(GetDefaultMongoDatabase(), &options.BucketOptions{
		Name: &collectionName,
	})
	if err != nil {
		return "", err
	}

	uploadOpts := options.GridFSUpload()

	r := bytes.NewBuffer(fileContent)

	if objectId, err := bucket.UploadFromStream(filename, r, uploadOpts); err != nil {
		return "", err
	} else {
		return objectId.Hex(), nil
	}
}

func GridFSOpenId(collectionName string, fileId primitive.ObjectID) (file *gridfs.DownloadStream, err error) {
	bucket, err := gridfs.NewBucket(GetDefaultMongoDatabase(), &options.BucketOptions{
		Name: &collectionName,
	})
	if err != nil {
		return nil, err
	}

	return bucket.OpenDownloadStream(fileId)
}

func GridFSOpen(collectionName, filePath string) (file *gridfs.DownloadStream, err error) {
	bucket, err := gridfs.NewBucket(GetDefaultMongoDatabase(), &options.BucketOptions{
		Name: &collectionName,
	})
	if err != nil {
		return nil, err
	}

	return bucket.OpenDownloadStreamByName(filePath)
}

func GridFSRename(collectionName, newFilename string, fileId primitive.ObjectID) error {
	bucket, err := gridfs.NewBucket(GetDefaultMongoDatabase(), &options.BucketOptions{
		Name: &collectionName,
	})
	if err != nil {
		return err
	}

	return bucket.Rename(fileId, newFilename)
}
