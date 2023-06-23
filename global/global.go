package g

import (
	_ "embed"

	"service/config"

	db "service/pkg/database"
	"service/pkg/logging"
	"service/pkg/translator"

	"github.com/kataras/iris/v12"
)

//go:embed version
var Version string

//go:embed name
var Name string

var (
	// Url
	UserKey      = "user"
	TranslateKey = "translate"

	// Context
	WriterLock   = "WriterLock"
	ClosedWriter = "ClosedWriter"

	// Regex
	UuidRegex string = `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`
)

// Config
var CFG *config.Config = nil

// Utilities
var Logger logging.Logger = nil
var Translator translator.Translator = nil

// App
var App *iris.Application = nil

// Default DB
var DB db.RelationalDatabaseFunction = nil

// Connections
var AllSQLCons = map[string]db.RelationalDatabaseFunction{}
