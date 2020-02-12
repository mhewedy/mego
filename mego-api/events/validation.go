package events

import (
	"github.com/mhewedy/httputil"
	"github.com/mhewedy/mego/commons"
	"net/http"
)

func parseAndValidateSearchInput(r *http.Request) (*searchInput, error) {
	var i searchInput
	err := commons.JSONDecode(r.Body, &i)
	if err != nil {
		return nil, err
	}
	if len(i.Emails) == 0 {
		return nil, httputil.NewClientError("empty emails")
	}
	if len(i.Rooms) == 0 {
		return nil, httputil.NewClientError("empty rooms")
	}

	return &i, nil
}

func parseAndValidateCreateInput(r *http.Request) (*createInput, error) {
	var i createInput
	err := commons.JSONDecode(r.Body, &i)
	if err != nil {
		return nil, err
	}
	if len(i.To) == 0 {
		return nil, httputil.NewClientError("empty emails")
	}
	if len(i.Room) == 0 {
		return nil, httputil.NewClientError("room should be supplied")
	}

	if i.Duration <= 0 {
		return nil, httputil.NewClientError("duration is invalid")
	}

	return &i, nil
}
