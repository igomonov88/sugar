package user_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/igomonov88/sugar/internal/platform/auth"
	"github.com/igomonov88/sugar/internal/tests"
	"github.com/igomonov88/sugar/internal/user"
	"github.com/pkg/errors"
	"testing"
	"time"
)

// TestUser validates the full set of CRUD operations on User values.
func TestUser(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	t.Log("Given the need to work with User records.")
	{
		ctx := tests.Context()
		now := time.Date(2020, time.February, 3, 0,0,0,0, time.UTC)

		claims := auth.NewClaims(uuid.New().String(), now, time.Hour)

		nu := user.NewUser{
			Name:            "Igor Gomonov",
			Email:           "gomonov.igor@gmail.com",
			Password:        "qwerty",
			PasswordConfirm: "qwerty",
		}

		u, err := user.Create(ctx, db, nu, now)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to create user : %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to create user.", tests.Success)

		savedU, err := user.Retrieve(ctx, claims, db, u.ID)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve user by ID: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve user by ID.", tests.Success)

		if diff := cmp.Diff(u, savedU); diff != "" {
			t.Fatalf("\t%s\tShould get back the same user. Diff:\n%s", tests.Failed, diff)
		}
		t.Logf("\t%s\tShould get back the same user.", tests.Success)

		upd := user.UpdateUser{
			Name:            tests.StringPointer("Maris Gomonov"),
			Email:           tests.StringPointer("maris@gmail.com"),
		}

		if err := user.Update(ctx, db, claims, u.ID, upd, now); err != nil {
			t.Fatalf("\t%s\tShould be able to update user : %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to update user.", tests.Success)

		savedU, err = user.Retrieve(ctx, claims, db, u.ID)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve user : %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve user.", tests.Success)

		if savedU.Name != *upd.Name {
			t.Errorf("\t%s\tShould be able to see updates to Name.", tests.Failed)
			t.Log("\t\tGot:", savedU.Name)
			t.Log("\t\tExp:", *upd.Name)
		} else {
			t.Logf("\t%s\tShould be able to see updates to Name.", tests.Success)
		}

		if savedU.Email != *upd.Email {
			t.Errorf("\t%s\tShould be able to see updates to Email.", tests.Failed)
			t.Log("\t\tGot:", savedU.Email)
			t.Log("\t\tExp:", *upd.Email)
		} else {
			t.Logf("\t%s\tShould be able to see updates to Email.", tests.Success)
		}

		if err := user.Delete(ctx, db, u.ID); err != nil {
			t.Fatalf("\t%s\tShould be able to delete user : %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to delete user.", tests.Success)

		savedU, err = user.Retrieve(ctx, claims, db, u.ID)
		if errors.Cause(err) != user.ErrNotFound {
			t.Fatalf("\t%s\tShould NOT be able to retrieve user : %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould NOT be able to retrieve user.", tests.Success)
	}
}

// TestAuthenticate validates the behavior around authenticating users.
func TestAuthenticate (t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	t.Log("Given the need to authenticate users")
	{
		t.Log("\tWhen handling a single User.")
		ctx := tests.Context()
		now := time.Date(2020, time.February, 3, 0, 0, 0, 0, time.UTC)

		nu := user.NewUser{
			Name:            "Tanya Gomonova",
			Email:           "tanya@gmail.com",
			Password:        "qwerty",
			PasswordConfirm: "qwerty",
		}

		u, err := user.Create(ctx, db, nu, now)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to create user : %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to create user.", tests.Success)

		claims, err := user.Authenticate(ctx, db, now, "tanya@gmail.com", "qwerty")
		if err != nil {
			t.Fatalf("\t%s\tShould be able to generate claims : %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to generate claims.", tests.Success)

		want := auth.Claims{}
		want.Subject = u.ID
		want.ExpiresAt = now.Add(time.Hour).Unix()
		want.IssuedAt = now.Unix()

		if diff := cmp.Diff(want, claims); diff != "" {
			t.Fatalf("\t%s\tShould get back the expected claims. Diff:\n%s", tests.Failed, diff)
		}
		t.Logf("\t%s\tShould get back the expected claims.", tests.Success)
	}
}