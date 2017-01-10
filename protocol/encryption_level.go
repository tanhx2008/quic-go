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
