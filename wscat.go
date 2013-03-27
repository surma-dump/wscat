package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"code.google.com/p/go.net/websocket"
	"github.com/voxelbrain/goptions"
	"github.com/voxelbrain/k"
)

var (
	options = struct {
		Protocols []string      `goptions:"-p, --protocol, description='Add protocol declaration in request'"`
		Origin    *url.URL      `goptions:"-o, --origin, description='Origin value for the request', obligatory"`
		Version   int           `goptions:"--websocket-version, description='Websocket version'"`
		Help      goptions.Help `goptions:"-h, --help, description='Show this help'"`
		Remainder goptions.Remainder
	}{
		Protocols: []string{},
		Version:   13,
	}
)

func init() {
	fs := goptions.NewFlagSet("wscat", &options)
	fs.ParseAndFail(os.Stderr, os.Args[1:])
	if len(options.Remainder) != 1 {
		fmt.Fprintf(os.Stderr, "Error: Need exactly one server address to connect to\n")
		fs.PrintHelp(os.Stderr)
		os.Exit(1)
	}

	_, err := url.Parse(options.Remainder[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Server address needs to be a valid URL\n")
		fs.PrintHelp(os.Stderr)
		os.Exit(1)
	}
}

func main() {
	cfg := &websocket.Config{
		Location: k.MustURL(options.Remainder[0]),
		Origin:   options.Origin,
		Protocol: options.Protocols,
		Version:  options.Version,
	}
	ws, err := websocket.DialConfig(cfg)
	if err != nil {
		log.Fatalf("Could not connect to server: %s", err)
	}
	go io.Copy(ws, os.Stdin)
	io.Copy(os.Stdout, ws)
}

type Header struct {
	Key, Value string
}

func (h *Header) MarshalGoptions(val string) error {
	fields := strings.SplitN(val, ":", 2)
	if len(fields) != 2 {
		return fmt.Errorf("Invalid header value")
	}

	h.Key = http.CanonicalHeaderKey(fields[0])
	h.Value = fields[1]
	return nil
}
