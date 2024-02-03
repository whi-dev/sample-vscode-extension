package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/sourcegraph/jsonrpc2"
)

type LanguageServer struct{}

func (s *LanguageServer) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	switch req.Method {
	case "textDocument/didOpen":
		params := struct {
			TextDocument struct {
				Text string `json:"text"`
			} `json:"textDocument"`
		}{}
		if err := req.Params(&params); err != nil {
			log.Printf("Error decoding textDocument/didOpen: %v", err)
			return
		}

		result := calculate(params.TextDocument.Text)
		conn.Notify(ctx, "calculatedResult", result)
	}
}

func calculate(input string) int {
	// 単純な計算を行う例
	parts := strings.Split(input, "+")
	if len(parts) != 2 {
		return 0
	}

	a, b := 0, 0
	fmt.Sscanf(parts[0], "%d", &a)
	fmt.Sscanf(parts[1], "%d", &b)

	return a + b
}

func main() {
	connOpt := jsonrpc2.ConnServerOption(jsonrpc2.MethodNotFound)
	conn := jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(stdrwc{}, jsonrpc2.VSCodeObjectCodec), &LanguageServer{}, connOpt)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn.Wait()
	}()
	wg.Wait()
}

// stdrwc is a simple io.ReadWriteCloser that delegates to os.Stdin and os.Stdout.
type stdrwc struct{}

func (stdrwc) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}

func (stdrwc) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (stdrwc) Close() error {
	return nil
}
