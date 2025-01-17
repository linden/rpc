package serverless

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/linden/rpc"
)

type Greeter struct{}

func (g *Greeter) Greet(req *string, res *string) error {
	fmt.Println("greet")

	*res = fmt.Sprintf("Hello %s", *res)

	return nil
}

func TestServerless(t *testing.T) {
	arith := new(Greeter)

	rpcsrv := rpc.NewServer()
	rpcsrv.Register(arith)

	h := NewHandler(rpcsrv)

	srv := httptest.NewServer(h)

	c := NewClient(srv.URL)

	a := "Jim"
	var b string

	err := c.Call("Greeter.Greet", &a, &b)
	if err != nil {
		t.Fatal(err)
	}

	a = "Jack"
	b = ""

	err = c.Call("Greeter.Greet", &a, &b)
	if err != nil {
		t.Fatal(err)
	}
}
