package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubclaunchpad/rocket/bot"
)

func TestPluginRegistration(t *testing.T) {
	b := bot.NewEmptyBot()
	err := RegisterPlugins(b)
	assert.Nil(t, err)
}
