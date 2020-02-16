package main

import (
	"encoding/json"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler-mail-outbound/sendgrid"
	sql "github.com/markdicksonjr/nibbler-sql"
	"github.com/markdicksonjr/nibbler-sql/session"
	userSql "github.com/markdicksonjr/nibbler-sql/user"
	"github.com/markdicksonjr/nibbler/build"
	"github.com/markdicksonjr/nibbler/session"
	"github.com/markdicksonjr/nibbler/user"
	"github.com/markdicksonjr/nibbler/user/auth/local"
	"net/http"
)

type Extension struct {
	nibbler.NoOpExtension
	AuthExtension *local.Extension
	UserExtension *user.Extension
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
	context.Router.HandleFunc("/api/ok", s.AuthExtension.SessionExtension.EnforceLoggedIn(func(w http.ResponseWriter, r *http.Request) {

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
		ctx.Count += 2

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

	// allocate our logger implementation that just uses Println
	logger := nibbler.DefaultLogger{}

	if build.GitTag != "" {
		logger.Info("running nibbler-sample v" + build.GitTag)
	}

	// allocate configuration, resolve env vars and config.json
	config, err := nibbler.LoadConfiguration()

	// any error is fatal at this point, this calls log.Fatal if non-nil
	nibbler.LogFatalNonNil(logger, err)

	// override the default static directory to where our vue app builds
	config.StaticDirectory = "./public/vue/dist"

	// allocate the sql extension, with all ORM-capable models - this extension auto-migrates the DB when initialized
	sqlExtension := sql.Extension{
		Models: []interface{}{
			nibbler.User{},
		},
	}

	// allocate user extension, providing sql extension to it - the user extension manages User CRUD
	userExtension := user.Extension{
		PersistenceExtension: &userSql.Extension{
			SqlExtension: &sqlExtension,
		},
	}

	// allocate session extension, using a sql connector - our sql connector will use our sql extension and tie into
	// the same DB as the operational database.  Allow SESSION_SECRET/SESSION_MAXAGE/SESSION_COOKIE env vars or
	// session.secret/session.maxage/session.cookie in the config.json file
	sessionExtension := session.Extension{
		StoreConnector: connectors.SqlStoreConnector{
			SqlExtension:  &sqlExtension,
			Secret:        config.Raw.Get("session", "secret").String("default_secret"),
			MaxAgeSeconds: config.Raw.Get("session", "maxage").Int(60 * 60 * 24 * 15), // 15 days
		},
		SessionName: config.Raw.Get("session", "cookie").String("default_cookie"),
	}

	// allocate the sendgrid extension
	sendgridExtension := sendgrid.Extension{}

	// allocate user local auth extension
	userLocalAuthExtension := local.Extension{
		SessionExtension:             &sessionExtension,
		UserExtension:                &userExtension,
		Sender:                       &sendgridExtension,
		PasswordResetEnabled:         config.Raw.Get("passwordreset", "enabled").Bool(true),
		PasswordResetFromName:        config.Raw.Get("passwordreset", "from", "name").String("Nibbler Sample"),
		PasswordResetFromEmail:       config.Raw.Get("passwordreset", "from", "email").String("noreply@nibblersample.com"),
		PasswordResetRedirect:        config.Raw.Get("passwordreset", "redirect").String("http://localhost:3000/#/reset-password"),
		RegistrationEnabled:          config.Raw.Get("registration", "enabled").Bool(true),
		RegistrationRequiresEmail:    config.Raw.Get("registration", "required", "email").Bool(true),
		RegistrationRequiresUsername: config.Raw.Get("registration", "required", "username").Bool(false),
		EmailVerificationEnabled:     config.Raw.Get("emailverification", "enabled").Bool(true),
		EmailVerificationFromName:    config.Raw.Get("emailverification", "from", "name").String("Nibbler Sample"),
		EmailVerificationFromEmail:   config.Raw.Get("emailverification", "from", "email").String("noreply@nibblersample.com"),
		EmailVerificationRedirect:    config.Raw.Get("emailverification", "redirect").String("http://localhost:3000/#/verify-email"),
	}

	// allocate our app-specific extension, providing the auth extension and group extension - note that the user
	// extension and others can be accessed through these extension instances
	sampleExtension := Extension{
		AuthExtension: &userLocalAuthExtension,
		UserExtension: &userExtension,
	}

	// initialize the application, which initializes all extensions in the order they are provided
	appContext := nibbler.Application{}
	nibbler.LogFatalNonNil(logger, appContext.Init(config, logger, []nibbler.Extension{
		&sqlExtension,
		&userExtension,
		&sessionExtension,
		&userLocalAuthExtension,
		&sendgridExtension,
		&sampleExtension,
	}))

	// create a test user, if it does not exist
	emailVal := "someone@example.com"
	usernameVal := "admin"
	password, _ := local.GeneratePasswordHash("tester123")
	_, _ = userExtension.Create(&nibbler.User{
		Email:    &emailVal,
		Username: &usernameVal,
		Password: &password,
	})

	// run the app
	nibbler.LogFatalNonNil(logger, appContext.Run())
}
