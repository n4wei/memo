package db

type Client interface {
	Get(string) ([]byte, bool)
	Set(string, []byte) bool
}
