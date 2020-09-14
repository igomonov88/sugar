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
		Description: "Add food table",
		Script: `
	CREATE TABLE IF NOT EXISTS food (
		id SERIAL PRIMARY KEY,
		fdc_id INT UNIQUE NOT NULL,
		description VARCHAR,
		brand_owner VARCHAR
	);`,
	},
	{
		Version:     2,
		Description: "Add search_food table",
		Script: `
	CREATE TABLE IF NOT EXISTS search_food (
		search_input VARCHAR NOT NULL,
		fdc_id INT NOT NULL
	);`,
	},
	{
		Version:     3,
		Description: "Add Carbohydrates table",
		Script: `
	CREATE TABLE IF NOT EXISTS carbohydrates (
		id SERIAL PRIMARY KEY, 
		fdc_id INT NOT NULL, 
		amount FLOAT NOT NULL, 
		unit_name VARCHAR,
		FOREIGN KEY (fdc_id) REFERENCES food(fdc_id));`,
	},
	{
		Version:     4,
		Description: "Add index to carbohydrates table",
		Script:      `CREATE INDEX idx_fdc_id on carbohydrates(fdc_id);`,
	},
	{
		Version:     5,
		Description: "Add portion table",
		Script: `
	CREATE TABLE IF NOT EXISTS portions (
		id SERIAL PRIMARY KEY,
		fdc_id INT NOT NULL,
		gram_weight FLOAT NOT NULL,
		description VARCHAR,
		FOREIGN KEY (fdc_id) REFERENCES food(fdc_id));
	`,
	},
}
