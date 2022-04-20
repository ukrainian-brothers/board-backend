package internal

import (
	"github.com/stretchr/testify/assert"
	"github.com/ukrainian-brothers/board-backend/internal/common"
	"testing"
)

func GetTestConfig(t *testing.T) *common.Config {
	cfg, err := common.NewConfigFromFile("../../config/configuration.test.local.json")
	assert.NoError(t, err)
	return cfg
}
