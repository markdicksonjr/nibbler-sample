package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler-sample/core"
	"github.com/markdicksonjr/nibbler/database/sql"
	"github.com/markdicksonjr/nibbler/mail/outbound/sendgrid"
	"github.com/markdicksonjr/nibbler/session"
	"github.com/markdicksonjr/nibbler/session/connectors"
	"github.com/markdicksonjr/nibbler/user"
	"github.com/markdicksonjr/nibbler/user/auth/local"
	userSql "github.com/markdicksonjr/nibbler/user/database/sql"
	"log"
)

func main() {

	// allocate configuration
	config, err := nibbler.LoadConfiguration(nil)

	// any error is fatal at this point
	if err != nil {
		log.Fatal(err.Error())
	}

	config.StaticDirectory = "./public/vue/dist"

	// allocate the sql extension, with all models
	sqlExtension := sql.Extension{
		Models: []interface{}{
			user.User{},
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

	// initialize the application
	appContext := nibbler.Application{}
	if err := appContext.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		&sqlExtension,
		&userExtension,
		&sessionExtension,
		&userLocalAuthExtension,
		&sendgridExtension,
		&core.Extension{
			AuthExtension: &userLocalAuthExtension,
		},
	}); err != nil {
		log.Fatal(err.Error())
	}

	// create a test user, if it does not exist
	emailVal := "someone@example.com"
	password, _ := local.GeneratePasswordHash("tester123")
	_, _ = userExtension.Create(&user.User{
		Email:    &emailVal,
		Password: &password,
	})

	// run the app
	if err = appContext.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
