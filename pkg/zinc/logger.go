package zinc

type Logger struct {
	Name string
}

func (l Logger) LogInfo(message string) {
	CreateDoc(l.Name, map[string]string{
		"severity": "info",
		"message":  message,
	})

}

func (l Logger) LogError(message, err string) {
	CreateDoc(l.Name, map[string]string{
		"severity": "error",
		"message":  message,
		"error":    err,
	})

}
