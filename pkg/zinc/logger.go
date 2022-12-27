package zinc

func LogInfo(index, message string) {
	CreateDoc(index, map[string]string{
		"severity": "info",
		"message":  message,
	})

}

func LogError(index, message, err string) {
	CreateDoc(index, map[string]string{
		"severity": "error",
		"message":  message,
		"error":    err,
	})

}
