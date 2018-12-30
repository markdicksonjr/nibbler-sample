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

	// allocate logger
	var logger nibbler.Logger = nibbler.DefaultLogger{}

	// allocate configuration
	config, err := nibbler.LoadConfiguration(nil)
	config.StaticDirectory = "./public/vue/dist"

	// any error is fatal at this point
	if err != nil {
		log.Fatal(err.Error())
	}

	// prepare models for initialization
	var models []interface{}
	models = append(models, user.User{})

	// allocate the sql extension, with all models
	sqlExtension := sql.Extension{
		Models: models,
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
			SqlExtension: &sqlExtension,
			Secret:      "dumbsecret",
			MaxAgeMs:    60 * 60 * 24 * 15, // 15 days
		},
		SessionName: "dumbcookie",
	}

	// allocate the sendgrid extension
	sendgridExtension := sendgrid.Extension{}

	// allocate user local auth extension
	userLocalAuthExtension := local.Extension{
		SessionExtension:       &sessionExtension,
		UserExtension:          &userExtension,
		Sender:     			&sendgridExtension,
		PasswordResetEnabled:   true,
		PasswordResetFromName:  "Nibbler Sample",
		PasswordResetFromEmail: "noreply@nibblersample.com",
		PasswordResetRedirect:  "http://localhost:3000/#/reset-password",
		RegistrationEnabled:	true,
		EmailVerificationEnabled:true,
		EmailVerificationFromName:  "Nibbler Sample",
		EmailVerificationFromEmail:"noreply@nibblersample.com",
		EmailVerificationRedirect:"http://localhost:3000/#/verify-email",
	}

	// prepare extensions for initialization
	extensions := []nibbler.Extension{
		&sqlExtension,
		&userExtension,
		&sessionExtension,
		&userLocalAuthExtension,
		&sendgridExtension,
		&core.Extension{
			AuthExtension: &userLocalAuthExtension,
		},
	}

	// initialize the application
	appContext := nibbler.Application{}
	if err := appContext.Init(config, &logger, &extensions); err != nil {
		log.Fatal(err.Error())
	}

	// run the app
	if err = appContext.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
