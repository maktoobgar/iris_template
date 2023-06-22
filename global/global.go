package g

import (
	_ "embed"
	"net/http"

	"service/config"

	"service/pkg/logging"
	"service/pkg/translator"
)

//go:embed version
var Version string

//go:embed name
var Name string

var (
	Query               = "query"
	UserKey             = "user"
	TranslateKey        = "translate"
	DeckId              = "deck_id"
	Deck                = "deck"
	Word                = "word"
	WordId              = "word_id"
	Type                = "type"
	UuidRegex    string = `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`
)

// Config
var CFG *config.Config = nil

// Utilities
var Logger logging.Logger = nil
var Translator translator.Translator = nil

// Microservices
var AuthMic *config.Microservice = nil
var DeckMic *config.Microservice = nil
var GameMic *config.Microservice = nil

// App
var Server *http.Server = nil
