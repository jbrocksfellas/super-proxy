package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	"github.com/things-go/go-socks5"
)

func main() {

	// created trafic store with max traffic 5 mb and max user traffic 3 mb
	maxUserTraffic := 2  // 2 mb
	maxTotalTraffic := 5 // 5 mb
	trafficStore := NewTrafficStore(int64(maxTotalTraffic), int64(maxUserTraffic))

	// Create a SOCKS5 server
	server := socks5.NewServer(
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
		socks5.WithAuthMethods([]socks5.Authenticator{socks5.UserPassAuthenticator{Credentials: credentials}}),
		socks5.WithConnectHandle(func(ctx context.Context, writer io.Writer, request *socks5.Request) error {
			noOfBytes := len(request.Bytes())
			username := request.AuthContext.Payload["username"]

			customerTraffic := trafficStore.getUserTraffic(username)
			if customerTraffic >= int64(maxUserTraffic) {
				return errors.New("Traffic limit exceeded for user: " + username)
			}

			totalTraffic := trafficStore.getTotalTraffic()
			if totalTraffic >= int64(maxTotalTraffic) {
				return errors.New("Traffic limit exceeded for all users.")
			}

			trafficStore.addTraffic(username, int64(noOfBytes))
			return nil
		}),
	)

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", "localhost:8000"); err != nil {
		panic(err)
	}
}
