package tftp

import (
	tftp "github.com/pin/tftp"
	"io"
	"os"
	"fmt"
	"time"
)

type TFTPServerService interface {
	Serve(done chan bool)
}

type TFTPServerWrapper struct {
	*tftp.Server
}

var _ TFTPServerService = &TFTPServerWrapper{}

func NewTFTPServer() (*TFTPServerWrapper, error) {
	srv := &TFTPServerWrapper{tftp.NewServer(readHandler, writeHandler)}
	return srv, nil
}

func (srv *TFTPServerWrapper) Serve(done chan bool)  {
	srv.Server.SetTimeout(60 * time.Second) // optional
	err := srv.Server.ListenAndServe(":6969") // blocks until srv.Server.Shutdown() is called
	if err != nil {
		fmt.Println("ListenAndServer returned an error.")
		fmt.Println(err)
	}
	done <- true
}



// readHandler is called when client starts file download from server
func readHandler(filename string, rf io.ReaderFrom) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	// Set transfer size before calling ReadFrom.
	//rf.(tftp.OutgoingTransfer).SetSize(myFileSize)
	n, err := rf.ReadFrom(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("%d bytes sent\n", n)
	return nil
}

// writeHandler is called when client starts file upload to server
func writeHandler(filename string, wt io.WriterTo) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	n, err := wt.WriteTo(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("%d bytes received\n", n)
	return nil
}
