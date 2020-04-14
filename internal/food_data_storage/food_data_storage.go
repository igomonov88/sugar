package food_data_storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func Search(ctx context.Context, db *sqlx.DB, searchInput string) ([]Food, error) {
	ctx, span := trace.StartSpan(ctx, "internal.food_data_storage.FoodSearch")
	defer span.End()

	var foods []Food
	const querySelectFood = `
	SELECT * FROM food WHERE EXISTS (
		SELECT fdc_id FROM search_food WHERE search_input LIKE '%' || $1 || '%'
	);`

	err := db.SelectContext(ctx, &foods, querySelectFood, searchInput)
	if err != nil {
		return nil, errors.Wrap(err, "selecting food")
	}

	return foods, nil
}

func AddFood(ctx context.Context, db *sqlx.DB, nf Food, foodSearchInput string) error {
	ctx, span := trace.StartSpan(ctx, "internal,food_data_storage.AddFood")
	defer span.End()

	const (
		queryAddFood       = `INSERT INTO food (fdc_id, description, brand_owner) VALUES ($1, $2, $3)`
		queryAddFoodSearch = `INSERT INTO search_food (search_input, fdc_id) VALUES ($1, $2)`
	)

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	_, err = tx.ExecContext(ctx, queryAddFood, &nf.FDCID, &nf.Description, &nf.BrandOwner)
	if err != nil {
		// TODO: handle rollback error
		tx.Rollback()
		return errors.Wrap(err, "inserting food to food")
	}
	_, err = tx.ExecContext(ctx, queryAddFoodSearch, foodSearchInput, &nf.FDCID)
	if err != nil {
		// TODO: handle rollback error
		tx.Rollback()
		return errors.Wrap(err, "inserting food search_food")
	}
	// TODO: handle commit error
	tx.Commit()
	return nil
}

func GetDetails(ctx context.Context, db *sqlx.DB, fdcID int) (*FoodDetails, error) {
	ctx, span := trace.StartSpan(ctx, "internal,food_data_storage.GetDetails")
	defer span.End()

	const (
		queryGetNutrients = `
		SELECT fn.type, fn.amount, n.number, n.name, n.rank, n.unit_name FROM food_nutrient AS fn 
			INNER JOIN nutrients AS n ON fn.nutrient_id=n.id WHERE fn.fdc_id=$1`

		queryGetDescription = `SELECT description from food where fdc_id=$1`
	)

	var fn []FoodNutrient
	var description string

	if err := db.SelectContext(ctx, &fn, queryGetNutrients, fdcID); err != nil {
		return nil, errors.Wrap(err, "selecting from food_nutrients")
	}

	if err := db.GetContext(ctx, &description, queryGetDescription, fdcID); err != nil {
		return nil, errors.Wrap(err, "selecting from food")
	}

	return &FoodDetails{Description: description, FoodNutrients: fn}, nil
}

func AddDetails(ctx context.Context, db *sqlx.DB, fdcID int, fd FoodDetails) error {
	ctx, span := trace.StartSpan(ctx, "internal,food_data_storage.AddDetails")
	defer span.End()

	const (
		queryAddFood          = `INSERT INTO food (fdc_id, description) VALUES ($1, $2) ON CONFLICT DO NOTHING;`
		queryAddNutrients     = `INSERT INTO nutrients (name, rank, unit_name, number) VALUES ($1, $2, $3, $4) RETURNING id;`
		queryAddFoodNutrients = `INSERT INTO food_nutrient (type, amount, nutrient_id, fdc_id) VALUES ($1,$2,$3,$4);`
	)

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	_, err = tx.ExecContext(ctx, queryAddFood, fdcID, fd.Description)
	if err != nil {
		// TODO: handle rollback error
		tx.Rollback()
		return errors.Wrap(err, "inserting food")
	}

	for i := 0; i <= len(fd.FoodNutrients)-1; i++ {
		var nID int
		row := tx.QueryRowContext(ctx, queryAddNutrients,
			fd.FoodNutrients[i].Name, fd.FoodNutrients[i].Rank, fd.FoodNutrients[i].UnitName, fd.FoodNutrients[i].Number)
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

		_, err = tx.ExecContext(ctx, queryAddFoodNutrients, fd.FoodNutrients[i].Type, fd.FoodNutrients[i].Amount, nID, fdcID)
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
