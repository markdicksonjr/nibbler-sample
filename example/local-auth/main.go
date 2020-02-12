package main

import (
	"encoding/json"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler-mail-outbound/sendgrid"
	sql "github.com/markdicksonjr/nibbler-sql"
	"github.com/markdicksonjr/nibbler-sql/session"
	userSql "github.com/markdicksonjr/nibbler-sql/user"
	"github.com/markdicksonjr/nibbler/session"
	"github.com/markdicksonjr/nibbler/user"
	"github.com/markdicksonjr/nibbler/user/auth/local"
	"log"
	"net/http"
)

type Extension struct {
	nibbler.NoOpExtension
	AuthExtension    *local.Extension
	UserExtension    *user.Extension
}

func (s *Extension) GetName() string {
	return "sample"
}

type AppSpecificContext struct {
	Count int
}

func (s *Extension) PostInit(context *nibbler.Application) error {

	// this adds the route "/api/ok" which increments a counter in protected context and a different counter in context
	// after hitting this endpoint, send a request to "/api/user" to ensure you don't see the protected context but
	// see the user-visible context.  Log out, log in, notice the result gets restored from the db
	context.Router.HandleFunc("/api/ok", s.AuthExtension.EnforceLoggedIn(func(w http.ResponseWriter, r *http.Request) {

		// get the user from the session
		caller, err := s.AuthExtension.SessionExtension.GetCaller(r)
		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}
		if caller == nil {
			nibbler.Write404Json(w)
			return
		}

		// get the protected context on the caller, increment a counter
		var protectedContext AppSpecificContext
		if caller.ProtectedContext != nil {
			if err := json.Unmarshal([]byte(*caller.ProtectedContext), &protectedContext); err != nil {
				nibbler.Write500Json(w, err.Error())
				return
			}
		}
		protectedContext.Count++

		// convert ctx to string, write it to protected context
		pc, _ := json.Marshal(protectedContext)
		if pc != nil {
			pcS := string(pc)
			caller.ProtectedContext = &pcS
		}

		// get the protected context on the caller, increment a counter
		var ctx AppSpecificContext
		if caller.Context != nil {
			if err := json.Unmarshal([]byte(*caller.Context), &ctx); err != nil {
				nibbler.Write500Json(w, err.Error())
				return
			}
		}
		ctx.Count+=2

		c, _ := json.Marshal(ctx)
		if c != nil {
			cS := string(c)
			caller.Context = &cS
		}

		// save the user back to the session
		if err := s.AuthExtension.SessionExtension.SetCaller(w, r, caller); err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}

		// save the user back to the database
		if err := s.UserExtension.Update(caller); err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "OK"}`))
	})).Methods("GET")
	return nil
}

func main() {

	// allocate configuration
	config, err := nibbler.LoadConfiguration()

	// any error is fatal at this point
	if err != nil {
		log.Fatal(err.Error())
	}

	config.StaticDirectory = "./public/vue/dist"

	// allocate the sql extension, with all models
	sqlExtension := sql.Extension{
		Models: []interface{}{
			nibbler.User{},
		},
	}

	// allocate user extension, providing sql extension to it
	userExtension := user.Extension{
		PersistenceExtension: &userSql.Extension{
			SqlExtension: &sqlExtension,
		},
	}

	// allocate session extension, using a sql connector
	// our sql connector will use our sql extension and
	// tie into the same DB as the operational database
	sessionExtension := session.Extension{
		StoreConnector: connectors.SqlStoreConnector{
			SqlExtension:  &sqlExtension,
			Secret:        "dumbsecret",
			MaxAgeSeconds: 60 * 60 * 24 * 15, // 15 days
		},
		SessionName: "dumbcookie",
	}

	// allocate the sendgrid extension
	sendgridExtension := sendgrid.Extension{}

	// allocate user local auth extension
	userLocalAuthExtension := local.Extension{
		SessionExtension:           &sessionExtension,
		UserExtension:              &userExtension,
		Sender:                     &sendgridExtension,
		PasswordResetEnabled:       true,
		PasswordResetFromName:      "Nibbler Sample",
		PasswordResetFromEmail:     "noreply@nibblersample.com",
		PasswordResetRedirect:      "http://localhost:3000/#/reset-password",
		RegistrationEnabled:        true,
		EmailVerificationEnabled:   true,
		EmailVerificationFromName:  "Nibbler Sample",
		EmailVerificationFromEmail: "noreply@nibblersample.com",
		EmailVerificationRedirect:  "http://localhost:3000/#/verify-email",
	}

	sampleExtension := Extension{
		AuthExtension:    &userLocalAuthExtension,
		UserExtension:    &userExtension,
	}

	// initialize the application
	appContext := nibbler.Application{}
	if err := appContext.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		&sqlExtension,
		&userExtension,
		&sessionExtension,
		&userLocalAuthExtension,
		&sendgridExtension,
		&sampleExtension,
	}); err != nil {
		log.Fatal(err.Error())
	}

	// create a test user, if it does not exist
	emailVal := "someone@example.com"
	password, _ := local.GeneratePasswordHash("tester123")
	_, _ = userExtension.Create(&nibbler.User{
		Email:    &emailVal,
		Password: &password,
	})

	// run the app
	if err = appContext.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
