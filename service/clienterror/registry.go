package clienterror

// Normalizer maps a preprocessed upstream error string to client-facing copy.
// Return ok=false to defer to the next rule.
type Normalizer func(preferChinese bool, raw string) (msg string, ok bool)

var normalizers []Normalizer

// Register appends a vendor/domain normalizer. Registration order is defined in normalize.go init().
func Register(fn Normalizer) {
	normalizers = append(normalizers, fn)
}

func runNormalizers(preferChinese bool, raw string) (string, bool) {
	for _, fn := range normalizers {
		if msg, ok := fn(preferChinese, raw); ok {
			return msg, true
		}
	}
	return "", false
}
