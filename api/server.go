package api

import (
	"fmt"
	"github.com/awilliams/couchdb-utils/util"
)

type Server struct {
	Couchdb string
	Uuid    string
	Vendor  struct {
		Name    string
		Version string
	}
	Version string
}

func (s *Server) path() string {
	return ""
}

func (s *Server) String() string {
	if s.Vendor.Name == "" {
		return fmt.Sprintf("v%s %s", s.Version, s.Couchdb)
	} else {
		return fmt.Sprintf("v%s - %s %s", s.Version, s.Vendor.Name, s.Couchdb)
	}
}

func (s Server) PP(printer util.Printer) {
	if s.Vendor.Name == "" {
		printer.Print("v%s\n%s", s.Version, s.Couchdb)
	} else {
		printer.Print("v%s - %s\n%s", s.Version, s.Vendor.Name, s.Couchdb)
	}
}

func (c Couchdb) GetServer() (Server, error) {
	server := new(Server)
	err := c.getJson(server, server.path())
	return *server, err
}
