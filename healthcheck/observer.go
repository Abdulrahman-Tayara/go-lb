package healthcheck

import "tayara/go-lb/models"

type IObserver interface {
	ServerDown(server *models.Server)

	ServerUp(server *models.Server)
}

type IObservable interface {
	Attach(observer IObserver)

	Detach(observer IObserver)

	Check()
}

type mockObserver struct {
	serversUpCount int

	serversDownCount int
}

func (m *mockObserver) ServerDown(server *models.Server) {
	m.serversDownCount++
}

func (m *mockObserver) ServerUp(server *models.Server) {
	m.serversUpCount++
}
