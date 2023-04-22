package main

import cmd "soul/cmd/app"

// @title						HandoverCloud Backend Server
// @version					v0.1
// @description				This is a Kubernetes-focused operations platform that includes CI/CD capabilities, gateway management, and service management features.
// @contact.name				SoulChild
// @contact.url				https://github.com/SoulChildTc/HandoverCloud
// @host						localhost:8080
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @BasePath					/
func main() {
	cmd.Execute()
}
