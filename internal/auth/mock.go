package auth

type mockStore struct {
	CheckCookieValueErr bool
}

func (m *mockStore) CheckCookieValue(value string) bool {
	return !m.CheckCookieValueErr
}

func (m *mockStore) InsertCookieValue(value string) (err error) {
	return nil
}
