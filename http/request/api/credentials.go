package api

type CredentialsForm struct {
	ServerId     string `json:"serverId" binding:"required"`
	ServerPasswd string `json:"serverPasswd" binding:"required"`
}
