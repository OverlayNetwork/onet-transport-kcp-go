package kcp

import (
	"context"
	"testing"

	"github.com/overlaynetwork/onet-go"
	"github.com/stretchr/testify/require"
)

func TestConn(t *testing.T) {

	laddr, err := onet.NewAddr("/ip/127.0.0.1/udp/1812/kcp")

	require.NoError(t, err)

	listener, err := onet.Listen(laddr)

	require.NoError(t, err)

	go func() {
		conn, err := onet.Dial(context.Background(), laddr)

		require.NoError(t, err)

		_, err = conn.Write([]byte("hello world"))

		require.NoError(t, err)
	}()

	_, err = listener.Accept()

	require.NoError(t, err)
}
