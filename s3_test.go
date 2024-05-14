package drive

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

var s3Client = &S3{
	Key:      os.Getenv("S3_KEY"),
	Secret:   os.Getenv("S3_SECRET"),
	Region:   os.Getenv("S3_REGION"),
	Endpoint: os.Getenv("S3_ENDPOINT"),
	Bucket:   os.Getenv("S3_BUCKET"),
}

const filePath = "testdata.txt"

func TestS3(t *testing.T) {
	// Read the testdata
	bytesToUpload, err := os.ReadFile(filePath)
	if err != nil {
		t.Error(err)
	}

	// Run the tests
	testS3Put(t, bytesToUpload)
	testS3Get(t, bytesToUpload)
	testS3Delete(t)
}

func testS3Put(t *testing.T, bytesToUpload []byte) {
	err := s3Client.Put(filePath, bytesToUpload)
	if err != nil {
		t.Error(err)
	}
}

func testS3Get(t *testing.T, uploadedBytes []byte) {
	// Retrieve an exising file from S3
	bytesFromGet, err := s3Client.Get(filePath)
	if err != nil {
		t.Error(err)
	}

	if string(uploadedBytes) != string(bytesFromGet) {
		t.Error("File content is not the same")
	}

	// Retrieve a non-existing file from S3
	randomFilePath := fmt.Sprintf("random-%d.txt", rand.Intn(1000))
	_, err = s3Client.Get(randomFilePath)
	if err == nil {
		t.Error("File should not exist")
	}
}

func testS3Delete(t *testing.T) {
	// Delete an existing file from S3
	err := s3Client.Delete(filePath)
	if err != nil {
		t.Error(err)
	}
}
