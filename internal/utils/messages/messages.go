package messages

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Get an i18nzed message (stored in <lang>.json files in assets folder)
func GetMessage(lang string, msg string) string {
	bundle := i18n.NewBundle(language.French)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	// load translation files
	bundle.LoadMessageFile("i18n/fr.toml")
	bundle.LoadMessageFile("i18n/en.toml")

	localizer := i18n.NewLocalizer(bundle, lang)

	localizeConfig := i18n.LocalizeConfig{
		MessageID: msg,
	}
	localization, _ := localizer.Localize(&localizeConfig)

	// fallback to initial message if no translation found
	if localization == "" {
		localization = msg
	}

	return localization
}
