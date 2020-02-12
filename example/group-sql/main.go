package main

import (
	"github.com/gorilla/mux"
	"github.com/markdicksonjr/nibbler"
	nibbler_sql "github.com/markdicksonjr/nibbler-sql"
	group_sql "github.com/markdicksonjr/nibbler-sql/group"
	connectors "github.com/markdicksonjr/nibbler-sql/session"
	user_sql "github.com/markdicksonjr/nibbler-sql/user"
	"github.com/markdicksonjr/nibbler/session"
	"github.com/markdicksonjr/nibbler/user"
	"github.com/markdicksonjr/nibbler/user/auth/local"
	"github.com/markdicksonjr/nibbler/user/group"
	"net/http"
)

type SampleExtension struct {
	nibbler.NoOpExtension
	UserGroupExtension *nibbler_user_group.Extension
	AuthExtension      *local.Extension
	AdminUser          *nibbler.User
}

func (s *SampleExtension) PostInit(app *nibbler.Application) error {
	app.Router.HandleFunc("/api/group", s.AuthExtension.EnforceLoggedIn(s.UserGroupExtension.CreateGroupRequestHandler)).Methods("PUT")
	app.Router.HandleFunc("/api/user/{id}/composite", s.AuthExtension.EnforceLoggedIn(s.UserGroupExtension.LoadUserCompositeRequestHandler)).Methods("GET", "POST")

	app.Router.HandleFunc("/api/auth", func(w http.ResponseWriter, r *http.Request) {
		if err := s.AuthExtension.SessionExtension.SetCaller(w, r, s.AdminUser); err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}
		nibbler.Write200Json(w, "{\"result\": \"ok\"}")
	}).Methods("GET")

	app.Router.HandleFunc("/api/noresource", s.UserGroupExtension.EnforceHasPrivilege("add", func(w http.ResponseWriter, r *http.Request) {
		nibbler.Write200Json(w, "{\"result\": \"ok\"}")
	})).Methods("GET")

	app.Router.HandleFunc("/api/{resource}", s.UserGroupExtension.EnforceHasPrivilegeOnResource("edit", func(r *http.Request) (s string, err error) {
		return mux.Vars(r)["resource"], nil
	}, func(w http.ResponseWriter, r *http.Request) {
		nibbler.Write200Json(w, "{\"result\": \"ok\"}")
	})).Methods("GET")

	return nil
}

func main() {
	logger := nibbler.DefaultLogger{}
	config, err := nibbler.LoadConfiguration()
	nibbler.LogFatalNonNil(logger, err)

	// allocate the sql extension, with all models
	sqlExtension := nibbler_sql.Extension{
		Models: nibbler_user_group.GetModels(),
	}

	// allocate user extension, providing sql extension to it
	userExtension := user.Extension{
		PersistenceExtension: &user_sql.Extension{
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

	// allocate our core extension
	coreExtension := nibbler_user_group.Extension{
		Logger: nibbler.DefaultLogger{},
		PersistenceExtension: &group_sql.SqlPersistenceExtension{
			SqlExtension: &sqlExtension,
		},
		SessionExtension: &sessionExtension,
		UserExtension:    &userExtension,
	}

	// allocate user local auth extension
	userLocalAuthExtension := local.Extension{
		SessionExtension: &sessionExtension,
		UserExtension:    &userExtension,
	}

	sampleExtension := SampleExtension{
		AuthExtension:      &userLocalAuthExtension,
		UserGroupExtension: &coreExtension,
	}

	// initialize the application
	appContext := nibbler.Application{}
	nibbler.LogFatalNonNil(logger, appContext.Init(config, logger, []nibbler.Extension{
		&sqlExtension,
		&userExtension,
		&sessionExtension,
		&userLocalAuthExtension,
		&coreExtension,
		&sampleExtension,
	}))

	// create a test admin user, if it does not exist
	emailVal := "admin@example.com"

	u, _ := userExtension.GetUserByEmail(emailVal)
	if u == nil {
		trueVal := true
		password, _ := local.GeneratePasswordHash("tester123")
		user, _ := userExtension.Create(&nibbler.User{
			Email:            &emailVal,
			Password:         &password,
			IsEmailValidated: &trueVal,
		})

		sampleExtension.AdminUser = user

		// create the admin group, add our admin, then give admin resource-agnostic "add" rights
		adminGroup, _ := coreExtension.CreateGroup("admins")
		_, err := coreExtension.SetGroupMembership(adminGroup.ID, user.ID, "admin")
		nibbler.LogFatalNonNil(logger, err)
		nibbler.LogFatalNonNil(logger, coreExtension.AddPrivilegeToGroups([]string{adminGroup.ID}, "", "add"))

		// give our admin group "edit" rights to a fictional "store"
		nibbler.LogFatalNonNil(logger, coreExtension.AddPrivilegeToGroups([]string{adminGroup.ID}, "store:1234", "edit"))

		// set the admin's current group as "admins"
		nibbler.LogFatalNonNil(logger, userExtension.Update(&nibbler.User{
			ID:             user.ID,
			CurrentGroupID: &adminGroup.ID,
		}))
	} else {
		sampleExtension.AdminUser = u
	}

	// run the app
	nibbler.LogFatalNonNil(logger, appContext.Run())
}
