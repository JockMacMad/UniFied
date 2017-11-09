package tftp

import (
	"github.com/pin/tftp"
	"os"
	"fmt"
	"bytes"
	"strconv"
)

func TFTPDownloadFromServer(host string, port int, path string) {
	var buffer bytes.Buffer
	buffer.WriteString(host)
	buffer.WriteString(":")
	buffer.WriteString(strconv.Itoa(port))
	host_port := buffer.String()
	c, err := tftp.NewClient(host_port)
	if err !=nil {

	}
	wt, err := c.Receive(path, "octet")
	if err !=nil {

	}
	file, err := os.Create(path)
	n, err := wt.WriteTo(file)
	fmt.Printf("%d bytes received\n", n)
}