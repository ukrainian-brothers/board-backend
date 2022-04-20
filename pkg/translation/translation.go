package translation

import (
	"encoding/json"
)

type LanguageTag string // ISO 639-1

type LanguageTags []LanguageTag

func (l LanguageTags) FromStrings(tags []string) LanguageTags {
	langTags := l
	for _, tag := range tags {
		langTags = append(langTags, LanguageTag(tag))
	}
	return langTags
}

func (l LanguageTags) Empty() bool {
	if l == nil {
		return true
	}
	if len(l) == 0 {
		return true
	}

	return l[0] == ""
}

const (
	English   LanguageTag = "en"
	Polish    LanguageTag = "pl"
	Ukrainian LanguageTag = "ua"
)

var supportedLanguages = []LanguageTag{
	English, Polish, Ukrainian,
}

type MultilingualString map[LanguageTag]string

func (s MultilingualString) Empty() bool {
	if len(s) == 0 {
		return true
	}

	for k, v := range s {
		if k != "" && v != "" {
			return false
		}
	}
	return true
}

func (s MultilingualString) MarshalJSON() ([]byte, error) {
	s.RemoveUnsupported()

	a := map[LanguageTag]string{}
	for k, v := range s {
		a[k] = v
	}

	return json.Marshal(a)
}

func (s MultilingualString) RemoveUnsupported() {
	for lang := range s {
		isSupported := false
		for _, supportedLang := range supportedLanguages {
			if supportedLang == lang {
				isSupported = true
				continue
			}
		}
		if !isSupported {
			delete(s, lang)
		}
	}
}
