package nagioseasier

import (
	"fmt"
	"net"
	"strings"
)

// QueryHandler is used to send queries to the query handler socket
type QueryHandler struct {
	Address string
}

// Create a QueryHandler
func Create(address string) (*QueryHandler, error) {
	if address == "" {
		address = "/var/lib/nagios/rw/nagios.qh"
	}

	qh := &QueryHandler{Address: address}
	return qh, nil
}

// Query performs a query using the QueryHandler
func (qh *QueryHandler) Query(command string) (string, error) {
	conn, err := net.Dial("unix", qh.Address)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(fmt.Sprintf("#nagioseasier %s\000", command)))
	if err != nil {
		return "", err
	}

	buf := make([]byte, 4096)
	output := ""
	for {
		n, err := conn.Read(buf[:])

		if err != nil {
			return scrub(output), err
		}

		if n == 0 {
			// connection closed by socket
			return scrub(output), nil
		}

		output += string(buf[0:n])
	}
}

func scrub(input string) (output string) {
	// get rid of pesky null chars
	output = strings.Replace(input, "\000", "", -1)

	// get rid of fake newlines, lol
	output = strings.Replace(output, "\\n", "", -1)

	// chomp off trailing newlines
	output = strings.Trim(output, "\n")

	// trim spaces
	output = strings.TrimSpace(output)

	return
}
