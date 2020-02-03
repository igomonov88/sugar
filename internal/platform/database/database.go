package database

// Config is the required properties to use the database.
type Config struct {
	User       string
	Password   string
	Host       string
	Name       string
	DisableTLS bool
}

