package rpc

import (
	netrpc "net/rpc"
)

type Client struct {
	cli *netrpc.Client
}

func NewClient(addr string) (*Client, error) {
	cli, err := netrpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli}, nil
}

func (c *Client) SayHello(name string) (string, error) {
	var resp HelloResp
	err := c.cli.Call("DemoService.SayHello", HelloReq{Name: name}, &resp)
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}
