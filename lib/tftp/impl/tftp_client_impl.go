package main

import (
	spi "bitbucket.org/ecosse-hosting/unified/lib/tftp/spi"
	"github.com/hashicorp/go-plugin"
	tftp "github.com/pin/tftp"
	"fmt"
	"bytes"
	"strconv"
	"os"
)

// Here is a real implementation of the TFTPClient
type TFTPClientImpl struct{
	client	*tftp.Client
}

// Connect to the TFTP Server
func (tc *TFTPClientImpl) Connect(host string, port int) string {
	var buffer bytes.Buffer
	var err error
	buffer.WriteString(host)
	buffer.WriteString(":")
	buffer.WriteString(strconv.Itoa(port))
	host_port := buffer.String()
	tc.client, err = tftp.NewClient(host_port)
	if err !=nil {

	}
	return "Hello!"
}

// Download a file from the TFTP Server to this 'client'
func (tc *TFTPClientImpl) DownloadFile(path string) string {
	wt, err := tc.client.Receive(path, "octet")
	if err !=nil {

	}
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error: \n", err)
	}
	n, err := wt.WriteTo(file)
	fmt.Printf("%d bytes received\n", n)
	return "Hello!"
}

// Upload a file from this TFTP 'client' to the remote TFTP Server.
func (tc *TFTPClientImpl) UploadFile(path string) string {
	return "Hello!"
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"tftpclient": &spi.TFTPClientPlugin{Impl: new(TFTPClientImpl)},
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
	vv := &TFTPClientImpl{}
	fmt.Println("Unified - TFTP Plugin - Connecting")
	vv.Connect("127.0.0.1", 69)
	//fmt.Println("Unified - TFTP Plugin - Downloading")
	//vv.DownloadFile("tftp_client_impl.go")
}