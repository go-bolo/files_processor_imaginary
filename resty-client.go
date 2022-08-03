package files_processor_imaginary

import (
	"fmt"
	"os"

	files_processor "github.com/go-catupiry/files/processor"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ClientConfiguration struct {
	URL string
}

// lib https://github.com/h2non/imaginary

type Client struct {
	HTTP *resty.Client
	Cfg  *ClientConfiguration
}

func NewClient(cfg *ClientConfiguration) *Client {
	c := Client{
		Cfg: cfg,
	}

	return &c
}

func (c *Client) Resize(sourcePath string, destPath string, opts files_processor.Options) error {
	url := c.Cfg.URL + "/resize"

	httpClient := resty.New()

	f, err := os.Open(sourcePath)
	if err != nil {
		return errors.Wrap(err, "error on open file from sourcePath")
	}
	defer f.Close()

	if _, ok := opts["type"]; !ok {
		opts["type"] = "webp"
	}

	if _, ok := opts["nocrop"]; !ok {
		opts["nocrop"] = "false"
	}

	id := uuid.New()

	res, err := httpClient.R().
		SetQueryParams(opts).
		SetFileReader("file", id.String(), f).
		SetContentLength(true).
		SetOutput(destPath).
		Post(url)

	// execution error
	if err != nil {
		return errors.Wrap(err, "error on access imaginary api")
	}
	// http error
	if res.IsError() {
		logrus.WithFields(logrus.Fields{
			"error":  fmt.Sprintf("%+v\n", err),
			"status": res.StatusCode(),
			"bory":   res.String(),
		}).Error("Response error", err)

		return errors.New(res.String())
	}

	return nil
}

// // TODO!

// crop - Same as /crop endpoint.
// smartcrop - Same as /smartcrop endpoint.
// enlarge - Same as /enlarge endpoint.
// extract - Same as /extract endpoint.
// rotate - Same as /rotate endpoint.
// autorotate - Same as /autorotate endpoint.
// flip - Same as /flip endpoint.
// flop - Same as /flop endpoint.
// thumbnail - Same as /thumbnail endpoint.
// zoom - Same as /zoom endpoint.
// convert - Same as /convert endpoint.
// watermark - Same as /watermark endpoint.
// watermarkimage - Same as /watermarkimage endpoint.
// blur - Same as /blur endpoint.
