package rpc

import "testing"

func TestDemoServiceSayHello(t *testing.T) {
	svc := &DemoService{}
	resp := &HelloResp{}
	if err := svc.SayHello(HelloReq{Name: "codex"}, resp); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.Message != "hello, codex" {
		t.Fatalf("unexpected message: %s", resp.Message)
	}
}
