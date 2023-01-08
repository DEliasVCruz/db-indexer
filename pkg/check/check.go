package check

import (
	"errors"
	"log"
)

func Error(check string, err error) {
	if err != nil {
		switch check {
		case "requestCreation":
			log.Fatalf("client: could not create request with %s\n", err)
		case "fileOpen":
			log.Fatalf("file: could not open given file with %s\n", err)
		case "requestAction":
			log.Fatalf("client: error making http request with %s\n", err)
		case "responseRead":
			log.Fatalf("client: could not read response body with %s\n", err)
		default:
			log.Fatalf("error: something went wrong with %s\n", err)
		}
	}

}

func SearchStatus(status int) error {

	switch status {
	case 500:
		log.Printf("server: could not find match with status %d", status)
		return errors.New("could not find match")
	case 200:
		return nil
	}
	
	return errors.New("index server error")
}
