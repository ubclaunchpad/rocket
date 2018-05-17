package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetCommand(t *testing.T) {
	ctx := getTestContext("@rocket set")
	b := getTestBot()
	res, _, err := b.commands["set"].Execute(ctx)
	t.Log(res)
	assert.Nil(t, err)
}
func TestBioFieldExists(t *testing.T) {
	b := getTestBot()
	res := b.commands["set"].Options["biography"].Key
	assert.Equal(t, res, "biography")
}
