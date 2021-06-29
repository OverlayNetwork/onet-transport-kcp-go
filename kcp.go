package kcp

import (
	"context"
	"net"
	"sync"

	"github.com/libs4go/errors"
	"github.com/overlaynetwork/onet-go"
	kcpgo "github.com/xtaci/kcp-go"
)

type kcpTransport struct {
	sync.RWMutex
	listeners map[string]net.Listener
}

func newKCPTransport() *kcpTransport {
	return &kcpTransport{
		listeners: make(map[string]net.Listener),
	}
}

func (transport *kcpTransport) String() string {
	return transport.Protocol()
}

func (transport *kcpTransport) Protocol() string {
	return "kcp"
}

func (transport *kcpTransport) listener(addr *onet.Addr) (net.Listener, error) {
	transport.Lock()
	defer transport.Unlock()

	listener, ok := transport.listeners[addr.String()]

	if !ok {
		tcpAddr, _, err := addr.ResolveNetAddr()

		if err != nil {
			return nil, errors.Wrap(err, "kcp transport listen on %s error", addr)
		}

		listener, err = kcpgo.Listen(tcpAddr.String())

		if err != nil {
			return nil, errors.Wrap(err, "kcp transport listen on %s error", addr)
		}

		transport.listeners[addr.String()] = listener
	}

	return listener, nil
}

func (transport *kcpTransport) Server(ctx context.Context, network *onet.OverlayNetwork, addr *onet.Addr, next onet.Next) (onet.Conn, error) {

	listener, err := transport.listener(addr)

	if err != nil {
		return nil, err
	}

	conn, err := listener.Accept()

	if err != nil {
		return nil, err
	}

	return onet.ToOnetConn(conn, network, addr)
}

func (transport *kcpTransport) Client(ctx context.Context, network *onet.OverlayNetwork, addr *onet.Addr, next onet.Next) (onet.Conn, error) {

	tcpAddr, _, err := addr.ResolveNetAddr()

	if err != nil {
		return nil, errors.Wrap(err, "kcp transport connect to %s error", addr)
	}

	conn, err := kcpgo.Dial(tcpAddr.String())

	if err != nil {
		return nil, errors.Wrap(err, "kcp transport conn to %s error", addr)
	}

	return onet.ToOnetConn(conn, network, addr)
}

func (transport *kcpTransport) Close(network *onet.OverlayNetwork, addr *onet.Addr, next onet.NextClose) error {
	transport.Lock()
	listener, ok := transport.listeners[addr.String()]

	if ok {
		delete(transport.listeners, addr.String())
		return listener.Close()
	}

	transport.Unlock()

	return nil
}

var protocol = &onet.Protocol{Name: "kcp", Native: true}

func init() {

	if err := onet.RegisterProtocol(protocol); err != nil {
		panic(err)
	}

	if err := onet.RegisterTransport(newKCPTransport()); err != nil {
		panic(err)
	}
}
