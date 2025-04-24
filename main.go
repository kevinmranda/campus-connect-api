package main

import (
	initializers "github.com/group4/campus-connect-api/Initializers"
	migrations "github.com/group4/campus-connect-api/Migrations"
	routes "github.com/group4/campus-connect-api/Routes"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	migrations.SyncDatabase()
	routes.Routes()
}

func main() {

 }