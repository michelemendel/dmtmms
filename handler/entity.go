package handler

type HandlerContext struct {
	loggedInUsers []LoggedInUser
}

func NewHandlerContext() *HandlerContext {
	return &HandlerContext{
		loggedInUsers: make([]LoggedInUser, 0),
	}
}
