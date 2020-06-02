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
	SELECT * FROM food WHERE fdc_id IN (
	    SELECT fdc_id FROM search_food WHERE search_input LIKE '%' || $1 ||'%');`

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
func RetrieveDetails(ctx context.Context, db *sqlx.DB, fdcID int) (*DetailsRef, error) {
	ctx, span := trace.StartSpan(ctx, "internal.storage.RetrieveDetailsRef")
	defer span.End()

	const (
		descriptionAndCarbsInfo = `
		SELECT f.description, c.amount, c.unit_name FROM food AS f 
		INNER JOIN carbohydrates AS c ON f.fdc_id = c.fdc_id and c.fdc_id = $1 
		FOR UPDATE;`
		portionsInfo = `
		SELECT id, fdc_id, gram_weight, description
		FROM portions WHERE fdc_id = $1;`
	)

	var portions []Portion

	err := db.SelectContext(ctx, &portions, portionsInfo, fdcID)
	if err != nil {
		return nil, err
	}

	details := DetailsRef{
		Portions: make([]Portion, len(portions)),
	}

	for i := range portions {
		details.Portions[i].ID = portions[i].ID
		details.Portions[i].FDCID = portions[i].FDCID
		details.Portions[i].Description = portions[i].Description
		details.Portions[i].GramWeight = portions[i].GramWeight
	}

	err = db.GetContext(ctx, &details, descriptionAndCarbsInfo, fdcID)
	if err != nil {
		return nil, err
	}

	return &details, nil
}

// SaveDetails save provided details to database.
func SaveDetails(ctx context.Context, db *sqlx.DB, fdcID int, carbs Carbohydrates, portions []Portion) error {
	ctx, span := trace.StartSpan(ctx, "internal.storage.SaveDetails")
	defer span.End()

	const (
		addCarbs = `INSERT INTO carbohydrates (fdc_id, amount, unit_name) 
		VALUES ($1, $2, $3) ON CONFLICT DO NOTHING;`
		addPortions = `INSERT INTO portions (fdc_id, gram_weight, description) 
		VALUES ($1, $2, $3) ON CONFLICT DO NOTHING;`
	)

	if len(portions) == 0 {
		_, err := db.Exec(addCarbs, fdcID, carbs.Amount, carbs.UnitName)
		if err != nil {
			return errors.Wrap(err, "inserting carbohydrates")
		}
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "creating transaction")
	}

	_, err = tx.Exec(addCarbs, fdcID, carbs.Amount, carbs.UnitName)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "inserting carbohydrates")
	}

	for i := range portions {
		_, err = tx.Exec(addPortions, fdcID,
			portions[i].GramWeight, portions[i].Description,
		)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "inserting portion")
		}
	}

	tx.Commit()
	return nil
}
