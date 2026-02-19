package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
)

func main() {
	stdScanner := bufio.NewScanner(os.Stdin)
	stdWriter := bufio.NewWriter(os.Stdout)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(fmt.Errorf("error connecting: %w", err))
	}
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			panic(fmt.Errorf("error close connecting: %w", err))
		}
	}(conn)

	srvReader := bufio.NewReader(conn)
	srvWriter := bufio.NewWriter(conn)

	for stdScanner.Scan() {
		input := stdScanner.Text()

		arr := strings.Split(input, " ")
		var output string

		if len(arr) > 0 {
			output = arr[0]
			if len(arr) > 1 {
				msg := strings.Join(arr[1:], " ")
				output = strings.Join(
					[]string{
						arr[0],
						base64.StdEncoding.EncodeToString([]byte(msg)),
					}, " ")
			}

		}

		_, err = srvWriter.WriteString(base64.StdEncoding.EncodeToString([]byte(output)) + "\n")
		if err != nil {
			slog.Error("Failed to write to server", slog.Any("error", err))
		}
		err = srvWriter.Flush()
		if err != nil {
			slog.Error("Failed to flush server", slog.Any("error", err))
		}

		r, err := srvReader.ReadString('\n')
		if err != nil {

			slog.Error("Failed to read from server", slog.Any("error", err))

		}

		_, err = stdWriter.WriteString(r)
		if err != nil {
			slog.Error("Failed to write to stdout", slog.Any("error", err))
		}
		err = stdWriter.Flush()
		if err != nil {
			slog.Error("Failed to flush stdout", slog.Any("error", err))
		}
		fmt.Printf("Output :%s\n", output)

	}
}
