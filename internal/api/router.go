package api

func (s *Server) SetupRouter() {
	s.router.All("/ping", s.ping)

}
