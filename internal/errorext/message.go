package errorext

var (
	ErrUserNotFound              = New("user does not exist")
	ErrNoActivePackage           = New("you don't have any active package")
	ErrUserBanned                = New("user banned")
	ErrUserOrPasswordIsIncorrect = New("user or password is incorrect")
)
