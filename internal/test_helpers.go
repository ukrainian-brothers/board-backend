package internal

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ukrainian-brothers/board-backend/internal/common"
	"testing"
)

func GetTestConfig(t *testing.T) *common.Config {
	cfg, err := common.NewConfigFromFile("../../config/configuration.test.local.json")
	assert.NoError(t, err)
	return cfg
}

var humanFriendlyUUIDMap = map[string]uuid.UUID{}

func HumanFriendlyUUID(s string) uuid.UUID {
	if humanFriendlyUUIDMap[s] == uuid.Nil {
		humanFriendlyUUIDMap[s] = uuid.New()
	}
	return humanFriendlyUUIDMap[s]
}
