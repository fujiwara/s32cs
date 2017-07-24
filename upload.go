package s32cs

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudsearchdomain"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

// const MaxUploadSize = 5 * 1024 * 1024
const MaxUploadSize = 500

var (
	openBracket  = []byte{'['}
	closeBracket = []byte{']'}
	comma        = []byte{','}
)

type buffer struct {
	bytes.Buffer
}

func (b *buffer) init() {
	b.Reset()
	b.Write(openBracket)
}

func (b *buffer) close() {
	b.Write(closeBracket)
}

func (b *buffer) allowAppend(bs []byte) bool {
	return b.Len()+len(bs)+2 < MaxUploadSize
}

func (b *buffer) append(bs []byte) {
	if b.Len() > 1 {
		b.Write(comma)
	}
	b.Write(bs)
}

type Domain struct {
	domain *cloudsearchdomain.CloudSearchDomain
	s3     *s3manager.Downloader
	buf    *buffer
}

func NewDomain(s1 *session.Session, s2 *session.Session) *Domain {
	buf := &buffer{}
	buf.Grow(MaxUploadSize)
	buf.init()
	return &Domain{
		domain: cloudsearchdomain.New(s1),
		s3:     s3manager.NewDownloader(s2),
		buf:    buf,
	}
}

func (d *Domain) Process(event S3Event) error {
	for _, record := range event.Records {
		r, err := d.fetch(record.S3.Bucket.Name, record.S3.Object.Key)
		if err != nil {
			return errors.Wrap(err, "fetch failed")
		}
		defer r.Close()

		if err = d.Upload(r); err != nil {
			return errors.Wrap(err, "upload failed")
		}
	}
	return nil
}

func (d *Domain) fetch(bucket, key string) (io.ReadCloser, error) {
	tmp, err := ioutil.TempFile(os.TempDir(), "s32cs")
	if err != nil {
		return nil, errors.Wrap(err, "create tempfile failed")
	}
	defer os.Remove(tmp.Name())

	log.Printf("downloading s3://%s/%s", bucket, key)
	n, err := d.s3.Download(tmp, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, errors.Wrap(err, "download failed")
	}
	log.Printf("%d bytes fetched", n)
	tmp.Seek(0, os.SEEK_SET)
	return tmp, nil
}

func (d *Domain) Upload(src io.Reader) error {
	dec := json.NewDecoder(src)
	for {
		var obj interface{}
		err := dec.Decode(&obj)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("decode json failed %s", err)
			continue
		}
		bs, err := json.Marshal(obj)
		if err != nil {
			return err
		}
		if !d.buf.allowAppend(bs) {
			err := d.flush()
			if err != nil {
				return err
			}
		}
		d.buf.append(bs)
	}
	return d.flush()
}

func (d *Domain) flush() error {
	defer d.buf.init()
	d.buf.close()
	log.Printf("starting upload %d bytes", d.buf.Len())
	log.Println(string(d.buf.Bytes()))
	out, err := d.domain.UploadDocuments(
		&cloudsearchdomain.UploadDocumentsInput{
			ContentType: aws.String("application/json"),
			Documents:   bytes.NewReader(d.buf.Bytes()),
		},
	)
	if err != nil {
		return errors.Wrap(err, "UploadDocuments failed")
	}
	log.Println("upload completed", out.String())
	return nil
}