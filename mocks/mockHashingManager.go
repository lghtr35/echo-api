package mocks

type MockHashingManager struct {
}

func NewMockHashingManager() *MockHashingManager {
	return &MockHashingManager{}
}

func (h *MockHashingManager) GetHash(s string) (string, error) {
	return s, nil
}

func (h *MockHashingManager) Verify(hashed string, new string) (bool, error) {
	res, err := h.GetHash(new)
	if err != nil {
		return false, err
	}
	return res == hashed, nil
}
