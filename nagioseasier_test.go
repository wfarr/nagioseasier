package nagioseasier

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestQuery(t *testing.T) {
	fp := filepath.Join(os.TempDir(), "nagioseasier.qh")
	listener, err := net.ListenUnix("unix", &net.UnixAddr{Net: "unix", Name: fp})
	defer listener.Close()

	go func() {
		conn, err := listener.AcceptUnix()
		defer conn.Close()

		if err != nil {
			t.Error(err)
		}

		conn.Write([]byte(fmt.Sprintf("loltests")))
	}()

	qh := Create(fp)

	resp, err := qh.Query("help")
	if err != nil {
		t.Error(err)
	}
	if resp != "loltests" {
		t.Error(resp)
	}
}
