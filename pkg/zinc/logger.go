package zinc

import "encoding/json"

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
