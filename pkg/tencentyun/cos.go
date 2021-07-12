package tencentyun

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// CosStorage cos storage
type CosStorage struct {
	secretID  string
	secretKey string
	region    string
	bucket    string
	rootPath  string
	Client    *cos.Client
	Timeout   time.Duration
}

// NewCosStorage new cos storage
func NewCosStorage(secretID, secretKey, region, bucket, rootPath string, timeout time.Duration) *CosStorage {
	u, _ := url.Parse("https://" + bucket + ".cos." + region + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Timeout: timeout,
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
	return &CosStorage{
		secretID:  secretID,
		secretKey: secretKey,
		region:    region,
		bucket:    bucket,
		rootPath:  rootPath,
		Client:    client,
		Timeout:   timeout,
	}
}

// Put put file from reader
func (c *CosStorage) Put(name string, r *io.Reader) error {
	dest := c.rootPath + "/" + name
	_, err := c.Client.Object.Put(context.Background(), dest, *r, nil)
	return err
}

// PutFile put file from local filename
func (c *CosStorage) PutFile(remoteFileName, localFileName string) error {
	dest := c.rootPath + "/" + remoteFileName
	_, err := c.Client.Object.PutFromFile(context.Background(), dest, localFileName, nil)
	return err
}

// Get get file
func (c *CosStorage) Get(name string) ([]byte, error) {
	dest := c.rootPath + "/" + name
	resp, err := c.Client.Object.Get(context.Background(), dest, nil)
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	// fmt.Printf("%s\n", string(bs))
	return bs, err
}
