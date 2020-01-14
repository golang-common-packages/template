package hash

// Storage interface for hash package
type Storage interface {
	FNV32(text string) uint32
	FNV32a(text string) uint32
	FNV64(text string) uint64
	FNV64a(text string) uint64
	MD5(text string) string
	SHA1(text string) string
	SHA256(text string) string
	SHA512(text string) string
}
