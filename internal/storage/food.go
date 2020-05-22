package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// List used for getting the list of Food items.
func List(ctx context.Context, db *sqlx.DB, searchInput string) ([]Food, error) {
	ctx, span := trace.StartSpan(ctx, "internal.storage.Search")
	defer span.End()

	var foods []Food

	const selectFood = `
	SELECT * FROM food WHERE EXISTS (
		SELECT fdc_id FROM search_food WHERE search_input LIKE '%' || $1 || '%'
	);`

	err := db.SelectContext(ctx, &foods, selectFood, searchInput)

	return foods, err
}

// SaveSearchInput is saved provided food item with associated search input.
func SaveSearchInput(ctx context.Context, db *sqlx.DB, food Food, input string) error {
	ctx, span := trace.StartSpan(ctx, "internal,storage.AddFood")
	defer span.End()

	const (
		addFood = `INSERT INTO food 
		(fdc_id, description, brand_owner) VALUES ($1, $2, $3)`

		addFoodSearch = `INSERT INTO search_food 
		(search_input, fdc_id) VALUES ($1, $2)`
	)

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	_, err = tx.Exec(addFood, food.FDCID, food.Description, food.BrandOwner)
	if err != nil {
		// TODO: handle rollback error
		tx.Rollback()
		return errors.Wrap(err, "inserting food to food")
	}
	_, err = tx.Exec(addFoodSearch, input, food.FDCID)
	if err != nil {
		// TODO: handle rollback error
		tx.Rollback()
		return errors.Wrap(err, "inserting food search_food")
	}
	// TODO: handle commit error
	tx.Commit()
	return nil
}

// RetrieveDetails returns Details and error if we got it.
func RetrieveDetails(ctx context.Context, db *sqlx.DB, fdcID int) (*Details, error) {
	ctx, span := trace.StartSpan(ctx, "internal,storage.GetDetails")
	defer span.End()

	var nutrients []FoodNutrient
	var description string

	const (
		getNutrients = `
		SELECT fn.type, fn.amount, n.number, n.name, n.rank, n.unit_name FROM 
        food_nutrient AS fn INNER JOIN nutrients AS n ON fn.nutrient_id=n.id 
        WHERE fn.fdc_id=$1`

		getDescription = `SELECT description from food where fdc_id=$1`
	)

	err := db.SelectContext(ctx, &nutrients, getNutrients, fdcID)
	if err != nil {
		return nil, errors.Wrap(err, "selecting from food_nutrients")
	}

	err = db.GetContext(ctx, &description, getDescription, fdcID)
	if err != nil {
		return nil, errors.Wrap(err, "selecting from food")
	}

	return &Details{Description: description, Nutrients: nutrients}, nil
}

// SaveDetails save providede details to database.
func SaveDetails(ctx context.Context, db *sqlx.DB, fdcID int, details Details) error {
	ctx, span := trace.StartSpan(ctx, "internal,storage.AddDetails")
	defer span.End()

	const (
		addFood = `INSERT INTO food 
		(fdc_id, description) VALUES ($1, $2) ON CONFLICT DO NOTHING;`

		addNutrients = `INSERT INTO nutrients 
		(name, rank, unit_name, number) VALUES ($1, $2, $3, $4) RETURNING id;`

		addFoodNutrients = `INSERT INTO food_nutrient 
		(type, amount, nutrient_id, fdc_id) VALUES ($1,$2,$3,$4);`
	)

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	_, err = tx.Exec(addFood, fdcID, details.Description)
	if err != nil {
		// TODO: handle rollback error
		tx.Rollback()
		return errors.Wrap(err, "inserting food")
	}

	for i := range details.Nutrients {
		var nID int
		row := tx.QueryRow(addNutrients, details.Nutrients[i].Name,
			details.Nutrients[i].Rank,
			details.Nutrients[i].UnitName,
			details.Nutrients[i].Number,
		)
		if err != nil {
			// TODO: handle rollback error
			tx.Rollback()
			return errors.Wrap(err, "inserting nutrients")
		}

		if err := row.Scan(&nID); err != nil {
			// TODO: handle rollback error
			tx.Rollback()
			return errors.Wrap(err, "scanning nutrientID")
		}

		_, err = tx.Exec(addFoodNutrients, details.Nutrients[i].Type,
			details.Nutrients[i].Amount, nID, fdcID)
		if err != nil {
			// TODO: handle rollback error
			tx.Rollback()
			return errors.Wrap(err, "inserting food_nutrient")
		}
	}
	// TODO: handle commit error
	tx.Commit()
	return nil
}
