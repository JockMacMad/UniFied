package spi

import "net/rpc"
import (
	api "bitbucket.org/ecosse-hosting/unified/lib/tftp/api"
	"github.com/hashicorp/go-plugin"
)

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a TFTPClientRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return TFTPClientRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type TFTPClientPlugin struct {
	// Impl Injection
	Impl api.TFTPClient
}

func (p *TFTPClientPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &api.TFTPClientRPCServer{Impl: p.Impl}, nil
}

func (TFTPClientPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &api.TFTPClientRPC{Client: c}, nil
}

