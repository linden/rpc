package serverless

import (
	"bytes"
	"io"
	"net/http"
	"sync"

	"github.com/linden/rpc"
)

var _ rpc.ClientCodec = (*clientGobCodec)(nil)

type clientGobCodec struct {
	url string
	res *http.Response

	gc   *rpc.GobClientCodec
	cond *sync.Cond
}

// Close implements rpc.ClientCodec.
func (c *clientGobCodec) Close() error {
	c.res = nil
	return nil
}

// ReadResponseBody implements rpc.ClientCodec.
func (c *clientGobCodec) ReadResponseBody(b any) error {
	return c.gc.ReadResponseBody(b)
}

// ReadResponseHeader implements rpc.ClientCodec.
func (c *clientGobCodec) ReadResponseHeader(r *rpc.Response) error {
	c.cond.Wait()

	return c.gc.ReadResponseHeader(r)
}

// WriteRequest implements rpc.ClientCodec.
func (c *clientGobCodec) WriteRequest(r *rpc.Request, body any) error {
	b := new(bytes.Buffer)

	gc := rpc.NewGobClientCodec(nil, b)
	gc.WriteRequest(r, body)

	req, err := http.NewRequest(http.MethodPost, c.url, b)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	rb, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	c.gc = rpc.NewGobClientCodec(bytes.NewBuffer(rb), nil)
	c.cond.Signal()

	return nil
}

func NewClient(url string) *rpc.Client {
	var m sync.Mutex
	m.Lock()

	return rpc.NewClientWithCodec(&clientGobCodec{
		url:  url,
		cond: sync.NewCond(&m),
	})
}
