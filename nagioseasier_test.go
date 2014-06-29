package nagioseasier

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func mockSocket() (*net.UnixListener, error) {
	fp := filepath.Join(os.TempDir(), "nagioseasier.qh")
	listener, err := net.ListenUnix("unix", &net.UnixAddr{Net: "unix", Name: fp})

	if err != nil {
		return nil, err
	}

	go func() {
		conn, _ := listener.AcceptUnix()
		defer conn.Close()
		conn.Write([]byte(fmt.Sprintf("loltests")))
	}()

	return listener, nil
}

func TestQuery(t *testing.T) {
	listener, err := mockSocket()
	if err != nil {
		t.Error(err)
	}

	defer listener.Close()

	qh := Create(listener.Addr().String())

	resp, err := qh.Query("help")
	if err != nil {
		t.Error(err)
	}
	if resp != "loltests" {
		t.Error(resp)
	}
}
