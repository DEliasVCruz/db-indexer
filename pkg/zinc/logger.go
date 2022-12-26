package zinc

func LogInfo(message string) {
	CreateDoc("appLogs", map[string]string{
		"severity": "info",
		"message":  message,
	})

}

func LogError(message, err string) {
	CreateDoc("appLogs", map[string]string{
		"severity": "error",
		"message":  message,
		"error":    err,
	})

}
