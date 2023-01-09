package check

import (
	"errors"
	"net/url"
)

var requiredParams = map[string][]string{
	"search": {"q", "from", "size", "field"},
}

func ParamsOf(endpoint string, params url.Values) error {

	for _, param := range requiredParams[endpoint] {
		if _, has := params[param]; !has {
			return errors.New("missing required params")
		}
	}

	return nil
}
