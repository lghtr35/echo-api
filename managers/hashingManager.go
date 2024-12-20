package managers

type HashingManager interface {
	GetHash(s string) (string, error)
	Verify(hashed string, new string) (bool, error)
}
