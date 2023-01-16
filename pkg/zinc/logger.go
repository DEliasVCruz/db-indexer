package zinc

import (
	"encoding/json"

	"github.com/DEliasVCruz/db-indexer/pkg/data"
)

func LogInfo(index, message string) {
	payLoad, _ := json.Marshal(
		map[string]string{
			"severity": "info",
			"message":  message,
		},
	)

	CreateDoc(index, payLoad)

}

func LogError(index, message, err string) {
	payLoad, _ := json.Marshal(
		map[string]string{
			"severity": "error",
			"message":  message,
			"error":    err,
		},
	)
	CreateDoc(index, payLoad)

}

func LogIndexStatus(upload bool, status, id string) error {
	response, err := json.Marshal(
		&data.FileUploaded{
			Uploaded: upload,
			State:    status,
			ID:       id,
		},
	)
	if err != nil {
		return err
	}

	CreateDoc("indexStatus", response)
	return nil

}
