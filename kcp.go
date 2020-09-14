package kcp

import (
	"context"

	"github.com/libs4go/errors"
	"github.com/overlaynetwork/onet-go"
	kcpgo "github.com/xtaci/kcp-go"
)

type kcpTransport struct{}

func (transport *kcpTransport) String() string {
	return transport.Protocol()
}

func (transport *kcpTransport) Protocol() string {
	return "kcp"
}

func (transport *kcpTransport) Listen(network *onet.OverlayNetwork) (onet.Listener, error) {

	tcpAddr, err := network.NavtiveAddr.ResolveNetAddr()

	if err != nil {
		return nil, errors.Wrap(err, "tcp transport listen on %s error", network.NavtiveAddr)
	}

	listen, err := kcpgo.Listen(tcpAddr.String())

	if err != nil {
		return nil, errors.Wrap(err, "tcp transport listen on %s error", network.NavtiveAddr)
	}

	return onet.ToOnetListener(listen, network)
}

func (transport *kcpTransport) Dial(ctx context.Context, network *onet.OverlayNetwork) (onet.Conn, error) {

	tcpAddr, err := network.NavtiveAddr.ResolveNetAddr()

	if err != nil {
		return nil, errors.Wrap(err, "tcp transport listen on %s error", network.NavtiveAddr)
	}

	conn, err := kcpgo.Dial(tcpAddr.String())

	if err != nil {
		return nil, errors.Wrap(err, "tcp transport conn to %s error", network.NavtiveAddr)
	}

	return onet.ToOnetConn(conn, network)
}

var protocol = &onet.Protocol{Name: "kcp"}

func init() {

	if err := onet.RegisterProtocol(protocol); err != nil {
		panic(err)
	}

	if err := onet.RegisterTransport(&kcpTransport{}); err != nil {
		panic(err)
	}
}
