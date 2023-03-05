package prompt

var DefaultLanguage = "English"

var languageMaps = map[string]string{
	"en":    DefaultLanguage,
	"zh-tw": "Traditional Chinese",
	"zh-cn": "Simplified Chinese",
	"ja":    "Japanese",
}

func GetLanguage(lang string) string {
	v, ok := languageMaps[lang]
	if !ok {
		return DefaultLanguage
	}
	return v
}
