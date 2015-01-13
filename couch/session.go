package couch

// GetUserCtx fetches the current user session, returning the UserCtx
func GetUserCtx(c *Client) (*UserCtx, error) {
	var result struct {
		OK      bool    `json:"ok"`
		UserCtx UserCtx `json:"userCtx"`
	}
	if err := c.Get("_session", nil, &result); err != nil {
		return nil, err
	}
	return &result.UserCtx, nil
}

// UserCtx is the part of the CouchDB user session information
type UserCtx struct {
	Name  string   `json:"name,omitempty"`
	Roles []string `json:"roles"`
}
