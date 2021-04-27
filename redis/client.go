package redis

import (
	"bufio"
	"io"
)

type Client struct {
	r *bufio.Reader
	w *bufio.Writer
}

func NewClient(stream io.ReadWriter) Client {
	return Client{
		r: bufio.NewReader(stream),
		w: bufio.NewWriter(stream),
	}
}

func (c Client) SetString(key, value string) (resp Response, err error) {
	err = c.sender().
		String(`SET "`).StringEscaped(key).String(`" "`).StringEscaped(value).String(`"`).
		Finish()
	if err != nil {
		return
	}
	return c.response(), nil
}

func (c Client) Exists(keys ...string) (resp Response, err error) {
	sender := c.sender().String("EXISTS")
	for _, k := range keys {
		sender.String(` "`).StringEscaped(k).String(`"`)
	}
	if err = sender.Finish(); err != nil {
		return
	}
	return c.response(), nil
}

func (c Client) Get(key string) (resp Response, err error) {
	err = c.sender().
		String(`GET "`).StringEscaped(key).String(`"`).
		Finish()
	if err != nil {
		return
	}
	return c.response(), nil
}

func (c Client) response() Response { return Response{r: c.r} }

func (c Client) sender() *sender { return newSender(c.w) }
