package test_helpers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/ukrainian-brothers/board-backend/internal/common"
	"math/rand"
	"testing"
)

func GetTestConfig(t *testing.T) *common.Config {
	cfg, err := common.NewConfigFromFile("../config/configuration.test.local.json")
	assert.NoError(t, err)
	return cfg
}

func RandomString(length int) string {
	by := make([]byte, length)
	for i := 0; i < length; i++ {
		by[i] = byte(65 + rand.Intn(25))
	}
	return string(by)
}

func RandomNumberRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func NewStringPtr(s string) *string {
	return &s
}

func RandomMail() *string {
	return NewStringPtr(fmt.Sprintf("%s@wp.pl", RandomString(8)))
}
