package prompt

const DefaultLanguage = "English"

var languageMaps = map[string]string{
	"en":    DefaultLanguage,
	"zh-tw": "Traditional Chinese",
	"zh-cn": "Simplified Chinese",
	"ja":    "Japanese",
	"pt":    "Portuguese",
	"pt-br": "Brazilian Portuguese",
}

// GetLanguage returns the language name for the given language code,
// or the default language if the code is not recognized.
func GetLanguage(langCode string) string {
	if language, ok := languageMaps[langCode]; ok {
		return language
	}
	return DefaultLanguage
}
