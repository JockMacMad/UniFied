package api

import "net/rpc"

// TFTPClient is the interface that we're exposing as a plugin.
type TFTPClient interface {
	Connect(host string, port int) string
	DownloadFile(path string) string
	UploadFile(path string) string
}

// Here is an implementation that talks over RPC
type TFTPClientRPC struct{ Client *rpc.Client }

type TFTPConnection struct {
	host	string
	port	int
}

func (g *TFTPClientRPC) Connect(host string, port int) string {
	var resp string
	err := g.Client.Call("Plugin.Connect", &TFTPConnection{host, port}, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

func (g *TFTPClientRPC) DownloadFile(path string) string {
	var resp string
	err := g.Client.Call("Plugin.DownloadFile", new(interface{}), &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

func (g *TFTPClientRPC) UploadFile(path string) string {
	var resp string
	err := g.Client.Call("Plugin.UploadFile", new(interface{}), &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Here is the RPC server that TFTPClientRPC talks to, conforming to
// the requirements of net/rpc
type TFTPClientRPCServer struct {
	// This is the real implementation
	Impl TFTPClient
}

func (s *TFTPClientRPCServer) Connect(args TFTPConnection, resp *string) error {
	*resp = s.Impl.Connect(args.host, args.port)
	return nil
}