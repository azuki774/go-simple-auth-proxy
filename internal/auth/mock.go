package auth

type mockStore struct {
	CheckCookieValueErr        bool
	ReturnGetBasicAuthPassword string
}

func (m *mockStore) CheckCookieValue(value string) bool {
	return !m.CheckCookieValueErr
}

func (m *mockStore) InsertCookieValue(value string) (err error) {
	return nil
}

func (m *mockStore) GetBasicAuthPassword(user string) string {
	return m.ReturnGetBasicAuthPassword
}
