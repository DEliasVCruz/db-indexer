package check

import (
	"errors"
	"net/url"
)

var requiredParams = map[string][]string{
	"search":      {"q", "from", "size", "field"},
	"indexStatus": {"id"},
}

func ParamsOf(endpoint string, params url.Values) error {

	for _, param := range requiredParams[endpoint] {
		if _, has := params[param]; !has {
			return errors.New("missing required params")
		}
	}

	return nil
}

func ValidPort(port int) error {
	if port < 1023 || port > 65535 {
		return errors.New("Invalid port number")
	}
	return nil
}
