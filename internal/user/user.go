package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/igomonov88/sugar/internal/platform/auth"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound is used when a specific User is requested but does not exist.
	ErrNotFound = errors.New("user not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("user_id is not on its proper form")

	// ErrAuthenticationFailure occurs when a user attempts to authenticate but
	// anything goes wrong.
	ErrAuthenticationFailure = errors.New("authentication failed")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")
)

// List retrieves a list of existing users from the database.
func List(ctx context.Context, db *sqlx.DB) ([]User, error) {
	ctx, span := trace.StartSpan(ctx, "internal.user.List")
	defer span.End()

	users := []User{}

	const q = `SELECT * FROM users;`

	if err := db.SelectContext(ctx, &users, q); err != nil {
		return nil, errors.Wrap(err, "selecting users")
	}
	return users, nil
}

// Retrieve gets the specified user from the database.
func Retrieve(ctx context.Context, claims auth.Claims, db *sqlx.DB, id string) (*User, error) {
	ctx, span := trace.StartSpan(ctx, "internal.user.Retrieve")
	defer span.End()

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrNotFound
	}

	var u User
	const q = `SELECT * FROM users WHERE user_id = $1;`
	if err := db.GetContext(ctx, &u, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
	}
	return &u, nil
}

// Create inserts a new user into the database.
func Create(ctx context.Context, db *sqlx.DB, nu NewUser, now time.Time) (*User, error) {
	ctx, span := trace.StartSpan(ctx, "internal.user.Create")
	defer span.End()

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "generating password hash")
	}

	u := User{
		ID:           uuid.New().String(),
		FirstName:    nu.FirstName,
		LastName:     nu.LastName,
		Email:        nu.Email,
		PasswordHash: hash,
		DateCreated:  now.UTC(),
		DateUpdated:  now.UTC(),
	}
	const q = `INSERT INTO users (user_id, first_name, last_name, email, password_hash, date_created, date_updated) 
				VALUES ($1, $2, $3, $4, $5, $6);`

	_, err = db.ExecContext(ctx, q, u.ID, u.FirstName, u.FirstName, u.Email, u.PasswordHash, u.DateCreated, u.DateUpdated)
	if err != nil {
		return nil, errors.Wrap(err, "inserting user")
	}
	return &u, nil
}

// Update replaces a user document in the database.
func Update(ctx context.Context, db *sqlx.DB, claims auth.Claims, id string, upd UpdateUser, now time.Time) error {
	ctx, span := trace.StartSpan(ctx, "internal.user.Update")
	defer span.End()

	u, err := Retrieve(ctx, claims, db, id)
	if err != nil {
		return err
	}

	if upd.FirstName != nil {
		u.FirstName = *upd.FirstName
	}

	if upd.LastName != nil {
		u.LastName = *upd.LastName
	}

	if upd.Email != nil {
		u.Email = *upd.Email
	}

	if upd.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*upd.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.Wrap(err, "generating password hash")
		}
		u.PasswordHash = hash
	}

	u.DateUpdated = now

	const q = `UPDATE users SET "first_name"=$2, "last_name"= $3, "email"=$4, "password_hash"=$5, "date_updated"=$6 WHERE user_id=$1`

	_, err = db.ExecContext(ctx, q, id, u.FirstName, u.LastName, u.Email, u.PasswordHash, u.DateUpdated)
	if err != nil {
		return errors.Wrap(err, "updating user")
	}
	return nil
}

// Delete removes a user from the database.
func Delete(ctx context.Context, db *sqlx.DB, id string) error {
	ctx, span := trace.StartSpan(ctx, "internal.user.Delete")
	defer span.End()

	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}

	const q = `DELETE FROM users WHERE user_id = $1`

	if _, err := db.ExecContext(ctx, q, id); err != nil {
		return errors.Wrapf(err, "deleting user %s", id)
	}

	return nil
}

// Authenticate finds a user by their email and verifies their password. On
// success it returns a Claims value representing this user. The claims can be
// used to generate a token for future authentication.
func Authenticate(ctx context.Context, db *sqlx.DB, now time.Time, email, password string) (auth.Claims, error) {
	ctx, span := trace.StartSpan(ctx, "internal.user.Authenticate")
	defer span.End()

	const q = `SELECT * FROM users WHERE email = $1`

	var u User
	if err := db.GetContext(ctx, &u, q, email); err != nil {
		if err == sql.ErrNoRows {
			return auth.Claims{}, ErrAuthenticationFailure
		}
		return auth.Claims{}, errors.Wrap(err, "selecting single user")
	}

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, ErrAuthenticationFailure
	}

	// If we are this far the request is valid. Create some claims for the user
	// and generate their token.
	claims := auth.NewClaims(u.ID, now, time.Hour)
	return claims, nil
}
