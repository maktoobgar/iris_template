package g

import (
	_ "embed"

	"service/config"

	db "service/pkg/database"
	"service/pkg/logging"
	media_manager "service/pkg/media"
	"service/pkg/translator"

	"github.com/kataras/iris/v12"
	"github.com/robfig/cron/v3"
)

//go:embed version
var Version string

//go:embed name
var Name string

var (
	// Header
	AccessToken = "Authorization"

	// Url
	TranslateKey = "translate"

	// Context
	WriterLock   = "WriterLock"
	ClosedWriter = "ClosedWriter"

	RequestBody = "RequestBody"
	DbInstance  = "DbInstance"
	UserKey     = "User"

	// Regex
	UuidRegex string = `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`
)

// Config
var CFG *config.Config = nil

// SecretKey in bytes
var SecretKeyBytes []byte

// Utilities
var Logger logging.Logger = nil
var Translator translator.Translator = nil

// App
var App *iris.Application = nil

// Main database type goes here
// Example: sqlite3 or postgres
var MainDatabaseType = ""

// Default DB
var DB db.RelationalDatabaseFunction = nil

// Connections
var AllSQLCons = map[string]db.RelationalDatabaseFunction{}

// Media manager for all medias
var Media media_manager.MediaManager = nil
var UsersMedia media_manager.MediaManager = nil

// Cron of the project
var Cron *cron.Cron = nil
