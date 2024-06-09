package auth

type mockStore struct {
	CheckCookieValueErr bool
}

func (m *mockStore) CheckCookieValue(value string) bool {
	return !m.CheckCookieValueErr
}
