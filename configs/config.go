package config

import "os"

var (
	DBHost = os.Getenv("DBHost")
	DBName = os.Getenv("DBName")
	DBUser = os.Getenv("DBUser")
	DBPwd  = os.Getenv("DBPwd")

	SessionKey = os.Getenv("session_key")
)
