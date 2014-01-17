package api

import (
	"github.com/awilliams/couchdb-utils/util"
	"strings"
)

type UserCtx struct {
	Name  string   `json:"name,omitempty"`
	Roles []string `json:"roles"`
}

type Session struct {
	OK      bool                   `json:"ok"`
	UserCtx UserCtx                `json:"userCtx"`
	Info    map[string]interface{} `json:"info"`
}

func (s Session) PP(printer util.Printer) {
	printer.Print("Name: %s\nRoles: %s", s.UserCtx.Name, strings.Join(s.UserCtx.Roles, ", "))
}

func (s *Session) path() string {
	return "_session"
}

func (c Couchdb) GetSession() (Session, error) {
	session := new(Session)
	err := c.getJson(session, session.path())
	return *session, err
}
