package nagioseasier

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

// QueryHandler is used to send queries to the query handler socket
type QueryHandler struct {
	Address *net.UnixAddr
}

// Create a QueryHandler
func Create(address string) *QueryHandler {
	if address == "" {
		address = "/var/lib/nagios/rw/nagios.qh"
	}

	addr := &net.UnixAddr{Name: address, Net: "unix"}
	qh := &QueryHandler{Address: addr}
	return qh
}

// Query performs a query with the QueryHandler socket
func (qh *QueryHandler) Query(command string) (string, error) {
	conn, err := net.DialUnix("unix", nil, qh.Address)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(fmt.Sprintf("#nagioseasier %s\000", command)))
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	n, err := buf.ReadFrom(conn)

	if n == 0 {
		return "", fmt.Errorf("no data received")
	}

	return scrub(buf.String()), err
}

func scrub(input string) string {
	var output string

	// get rid of pesky null chars
	output = strings.Replace(input, "\000", "", -1)

	// get rid of fake newlines, lol
	output = strings.Replace(output, "\\n", "", -1)

	// chomp off trailing newlines
	output = strings.Trim(output, "\n")

	// trim spaces
	output = strings.TrimSpace(output)

	return output
}
