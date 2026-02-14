package rpc

import "fmt"

type DemoService struct{}

func (s *DemoService) SayHello(req HelloReq, resp *HelloResp) error {
	name := req.Name
	if name == "" {
		name = "world"
	}
	resp.Message = fmt.Sprintf("hello, %s", name)
	return nil
}
