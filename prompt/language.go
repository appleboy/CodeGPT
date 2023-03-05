package prompt

var defaultLanguage = "English"

var languageMaps = map[string]string{
	"en":    defaultLanguage,
	"zh-tw": "Traditional Chinese",
	"zh-cn": "Simplified Chinese",
	"ja":    "Japanese",
}

func GetLanguage(lang string) string {
	v, ok := languageMaps[lang]
	if !ok {
		return defaultLanguage
	}
	return v
}
