package s32cs

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudsearchdomain"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

const MaxUploadSize = 5 * 1024 * 1024

var (
	openBracket  = []byte{'['}
	closeBracket = []byte{']'}
	comma        = []byte{','}
	DEBUG        = false
)

type Buffer struct {
	bytes.Buffer
}

func NewBuffer() *Buffer {
	buf := new(Buffer)
	buf.Grow(MaxUploadSize)
	buf.Init()
	return buf
}

func (b *Buffer) Init() {
	b.Reset()
	b.Write(openBracket)
}

func (b *Buffer) Close() {
	b.Write(closeBracket)
}

func (b *Buffer) allowAppend(bs []byte) bool {
	return b.Len()+len(bs)+2 < MaxUploadSize
}

func (b *Buffer) append(bs []byte) {
	if b.Len() > 1 {
		b.Write(comma)
	}
	b.Write(bs)
}

type Flusher func(*Buffer) error

type Client struct {
	endpoint string
	s3       *s3manager.Downloader
	queue    *sqs.SQS
	buf      *Buffer
	reg      *regexp.Regexp
}

func NewClient(sess *session.Session, endpoint string, reg *regexp.Regexp) *Client {
	return &Client{
		endpoint: endpoint,
		s3:       s3manager.NewDownloader(sess),
		queue:    sqs.New(sess),
		buf:      NewBuffer(),
		reg:      reg,
	}
}

func (c *Client) Process(event S3Event) error {
	for _, record := range event.Records {
		name, key := record.S3.Bucket.Name, record.S3.Object.Key
		r, err := c.fetch(name, key)
		if err != nil {
			return errors.Wrap(err, "fetch failed")
		}
		defer r.Close()

		endpoint := c.endpoint
		// extract endpoint from key
		if c.reg != nil {
			m := c.reg.FindStringSubmatch(key)
			switch len(m) {
			case 0:
				log.Printf("warn\textract endpoint from key %s by regexp %s failed. using default endpoint %s", key, c.reg.String(), endpoint)
			case 1:
				endpoint = m[0]
			default:
				endpoint = m[1]
			}
		}

		if err = c.Upload(r, endpoint); err != nil {
			return errors.Wrap(err, "[error] upload failed")
		}
	}
	return nil
}

func (c *Client) fetch(bucket, key string) (io.ReadCloser, error) {
	tmp, err := ioutil.TempFile(os.TempDir(), "s32cs")
	if err != nil {
		return nil, errors.Wrap(err, "create tempfile failed")
	}
	defer os.Remove(tmp.Name())

	log.Printf("info\tdownloading s3://%s/%s", bucket, key)
	n, err := c.s3.Download(tmp, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, errors.Wrap(err, "download failed")
	}
	log.Printf("info\t%d bytes fetched", n)
	tmp.Seek(0, os.SEEK_SET)

	if strings.HasSuffix(key, ".gz") {
		return gzip.NewReader(tmp)
	}
	return tmp, nil
}

func (c *Client) Upload(src io.Reader, endpoint string) error {
	if endpoint == "" {
		endpoint = c.endpoint
	}
	log.Printf("info\tendpoint %s", endpoint)
	f := func(buf *Buffer) error {
		return uploadCloudSearch(buf, endpoint)
	}
	return c.BuildAndFlush(src, f)
}

func (c *Client) BuildAndFlush(src io.Reader, flush Flusher) error {
	dec := json.NewDecoder(src)
	for {
		var record SDFRecord
		if err := dec.Decode(&record); err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("warn\tdecode json failed %s", err)
			continue
		}
		if err := record.Validate(); err != nil {
			log.Printf("warn\tSDF record validation failed %s %#v", err, record)
			continue
		}
		bs, err := json.Marshal(record)
		if err != nil {
			return err
		}
		if !c.buf.allowAppend(bs) {
			err := flush(c.buf)
			if err != nil {
				return err
			}
		}
		c.buf.append(bs)
	}
	return flush(c.buf)
}

func uploadCloudSearch(buf *Buffer, endpoint string) error {
	defer buf.Init()
	buf.Close()
	log.Printf("info\tstarting upload %d bytes", buf.Len())
	if DEBUG {
		log.Println("debug\t" + string(buf.Bytes()))
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String(endpoint),
	}))
	domain := cloudsearchdomain.New(sess)

	out, err := domain.UploadDocuments(
		&cloudsearchdomain.UploadDocumentsInput{
			ContentType: aws.String("application/json"),
			Documents:   bytes.NewReader(buf.Bytes()),
		},
	)
	if err != nil {
		return errors.Wrap(err, "UploadDocuments failed")
	}
	log.Println("info\tupload completed", out.String())
	return nil
}
