package quic

/*
import (
	"crypto/tls"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go/utils"
)

type ConnState int8

const (
	ConnStateVersionNegotiated ConnState = iota
	ConnStateEncryptedIncoming
	ConnStateEncryptedOutgoing
	ConnStateForwardSecureIncoming
	ConnStateForwardSecureOutgoing
)

type QuicListener struct {
	TLSConfig *tls.Config
	// called in a separate go routine
	ConnState func(*Session, ConnState)
	// Defaults to ...
	HandshakeTimeout time.Duration
	// more timeouts...
	//
	MaxIncomingStreams int

	// privates

	conn      net.PacketConn
	populated bool
}

func (l *QuicListener) ListenAddr(addr string) error     {}
func (l *QuicListener) Listen(conn net.PacketConn) error {}

// TODO: think about whether Dial is blocking (until when?) or non-blocking
func (l *QuicListener) DialAddr(addr string) error     {}
func (l *QuicListener) Dial(conn net.PacketConn) error {}

func (l *QuicListener) Close() error   {}
func (l *QuicListener) Addr() net.Addr {}

type SessionI interface {
	AcceptStream() (utils.Stream, error)
	// guaranteed to return the smallest unopened stream
	// special error for "too many streams, retry later"
	OpenStream() (utils.Stream, error)
	// blocks
	OpenStreamSync() (utils.Stream, error)

	RemoteAddr() net.Addr
	Close(error) error
}
*/
