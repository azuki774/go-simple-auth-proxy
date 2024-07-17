package auth

type mockStore struct {
	ReturnGetBasicAuthPassword string
}

func (m *mockStore) GetBasicAuthPassword(user string) string {
	return m.ReturnGetBasicAuthPassword
}
