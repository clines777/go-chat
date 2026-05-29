package protocol

import "fmt"

const (
	tokenPrefix = "token"
	loginToken  = "t_login"
	apiToken    = "t_api"
	resumeToken = "t_resume"
	session     = "session"
)

func LoginTokenKey(token string) string {
	return tokenPrefix + ":" + loginToken + ":" + token
}

func ApiTokenKey(token string) string {
	return tokenPrefix + ":" + apiToken + ":" + token
}

func ResumeTokenKey(token string) string {
	return tokenPrefix + ":" + resumeToken + ":" + token
}

func SessionKey(userID int, connID string) string {
	return fmt.Sprintf("%s:%d:%s", session, userID, connID)
}
