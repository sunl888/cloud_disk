package hasher

type Hasher interface {
	// Hash the given value.
	Make(value string) string

	// Check the given plain value against a hash.
	Check(value string, hashedValue string) bool
}
