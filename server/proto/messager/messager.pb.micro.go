// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: messager.proto

package messager

import (
	fmt "fmt"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	proto "google.golang.org/protobuf/proto"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Messager service

func NewMessagerEndpoints() []*api.Endpoint {
	return []*api.Endpoint{
		{
			Name:    "Messager.AddResource",
			Path:    []string{"/pumper/resource"},
			Method:  []string{"POST"},
			Handler: "rpc",
		},
		{
			Name:    "Messager.DeleteResource",
			Path:    []string{"/pumper/resource"},
			Method:  []string{"DELETE"},
			Handler: "rpc",
		},
	}
}

// Client API for Messager service

type MessagerService interface {
	AddResource(ctx context.Context, in *TaskSpec, opts ...client.CallOption) (*NodeSpec, error)
	DeleteResource(ctx context.Context, in *TaskSpec, opts ...client.CallOption) (*emptypb.Empty, error)
}

type messagerService struct {
	c    client.Client
	name string
}

func NewMessagerService(name string, c client.Client) MessagerService {
	return &messagerService{
		c:    c,
		name: name,
	}
}

func (c *messagerService) AddResource(ctx context.Context, in *TaskSpec, opts ...client.CallOption) (*NodeSpec, error) {
	req := c.c.NewRequest(c.name, "Messager.AddResource", in)
	out := new(NodeSpec)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagerService) DeleteResource(ctx context.Context, in *TaskSpec, opts ...client.CallOption) (*emptypb.Empty, error) {
	req := c.c.NewRequest(c.name, "Messager.DeleteResource", in)
	out := new(emptypb.Empty)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Messager service

type MessagerHandler interface {
	AddResource(context.Context, *TaskSpec, *NodeSpec) error
	DeleteResource(context.Context, *TaskSpec, *emptypb.Empty) error
}

func RegisterMessagerHandler(s server.Server, hdlr MessagerHandler, opts ...server.HandlerOption) error {
	type messager interface {
		AddResource(ctx context.Context, in *TaskSpec, out *NodeSpec) error
		DeleteResource(ctx context.Context, in *TaskSpec, out *emptypb.Empty) error
	}
	type Messager struct {
		messager
	}
	h := &messagerHandler{hdlr}
	opts = append(opts, api.WithEndpoint(&api.Endpoint{
		Name:    "Messager.AddResource",
		Path:    []string{"/pumper/resource"},
		Method:  []string{"POST"},
		Handler: "rpc",
	}))
	opts = append(opts, api.WithEndpoint(&api.Endpoint{
		Name:    "Messager.DeleteResource",
		Path:    []string{"/pumper/resource"},
		Method:  []string{"DELETE"},
		Handler: "rpc",
	}))
	return s.Handle(s.NewHandler(&Messager{h}, opts...))
}

type messagerHandler struct {
	MessagerHandler
}

func (h *messagerHandler) AddResource(ctx context.Context, in *TaskSpec, out *NodeSpec) error {
	return h.MessagerHandler.AddResource(ctx, in, out)
}

func (h *messagerHandler) DeleteResource(ctx context.Context, in *TaskSpec, out *emptypb.Empty) error {
	return h.MessagerHandler.DeleteResource(ctx, in, out)
}
