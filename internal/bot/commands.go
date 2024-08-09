package bot

var (
	MainBotCommands = map[string]string{
		"/start":       "show help message",
		"/status":      "show your active package status",
		"/create":      "create a user and show credentials, and activate a trial package if trial is activated by administrators",
		"/password":    "change your account password",
		"/connections": "show active connections",
	}
)
