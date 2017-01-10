package protocol

// EncryptionLevel is the encryption level
type EncryptionLevel int

const (
	// EncryptionUnspecified is not specified
	EncryptionUnspecified EncryptionLevel = iota
	// EncryptionUnencrypted is not encrypted
	EncryptionUnencrypted
	// EncryptionSecure is encrypted, but not forward secure
	EncryptionSecure
	// EncryptionForwardSecure is forward secure
	EncryptionForwardSecure
)

func (e EncryptionLevel) String() string {
	switch e {
	case EncryptionUnencrypted:
		return "unencrypted"
	case EncryptionSecure:
		return "encrypted"
	case EncryptionForwardSecure:
		return "forward encrypted"
	}
	return "unknown"
}
