package schema

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
)

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}

// migrations contains the queries needed to construct the database schema.
// Entries should never be removed from this slice once they have been ran in
// production.
//
// Using constants in a .go file is an easy way to ensure the queries are part
// of the compiled executable and avoids pathing issues with the working
// directory. It has the downside that it lacks syntax highlighting and may be
// harder to read for some cases compared to using .sql files. You may also
// consider a combined approach using a tool like packr or go-bindata.
var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add users",
		Script: `
CREATE TABLE users (
	user_id       UUID,
	name          TEXT,
	email         TEXT UNIQUE,
	password_hash TEXT,

	date_created TIMESTAMP,
	date_updated TIMESTAMP,

	PRIMARY KEY (user_id)
);`,
	},
	{
		Version:     2,
		Description: "Rename column name to first_name",
		Script:      `ALTER TABLE users RENAME COLUMN name TO first_name;`,
	},
	{
		Version:     3,
		Description: "add column last_name to users table",
		Script:      `ALTER TABLE users ADD COLUMN last_name TEXT;`,
	},
	{
		Version:     4,
		Description: "add roles field to users table",
		Script:      `ALTER TABLE users ADD COLUMN roles TEXT[]`,
	},
}
