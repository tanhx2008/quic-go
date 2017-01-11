package handshake

import "github.com/lucas-clemente/quic-go/protocol"

// CryptoSetup is a crypto setup
type CryptoSetup interface {
	Open(dst, src []byte, packetNumber protocol.PacketNumber, associatedData []byte) ([]byte, protocol.EncryptionLevel, error)

	// forceEncryptionLevel forces the packet to be encrypted with a certain encryption level. This is needed for the retransmission of handshake packets
	Seal(dst, src []byte, packetNumber protocol.PacketNumber, associatedData []byte, forceEncryptionLevel protocol.EncryptionLevel) ([]byte, protocol.EncryptionLevel, error)

	HandleCryptoStream() error
	DiversificationNonce() []byte
	LockForSealing()
	UnlockForSealing()
	HandshakeComplete() bool
}
