package translation

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslation(t *testing.T) {
	title := MultilingualString{
		"pl":                       "tytuł",
		"ua":                       "титул",
		"en":                       "title",
		"unsupported_language_tag": "x",
	}

	title.RemoveUnsupported()

	assert.Equal(t, "title", title[English])

	_, ok := title["unsupported_language_tag"]
	assert.Equal(t, false, ok)
}

func TestTranslationMarshal(t *testing.T) {
	title := MultilingualString{
		"ua":                       "титул",
		"unsupported_language_tag": "X",
	}

	_, err := json.Marshal(title)
	assert.Nil(t, err)
}

func TestTranslationEmpty(t *testing.T) {
	field := MultilingualString{}
	assert.Equal(t, true, field.Empty())
	field[Ukrainian] = ""
	assert.Equal(t, true, field.Empty())
	field[Ukrainian] = "x"
	assert.Equal(t, false, field.Empty())
}
