package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubclaunchpad/rocket/model"
)

func TestMemberCreateGetUpdateRemove(t *testing.T) {
	dal := newTestDBConnection()
	defer dal.Close()

	// Create a new member
	member := &model.Member{
		SlackID: "1234",
		Name:    "Big Bruno",
	}
	err := dal.CreateMember(member)
	assert.Nil(t, err)

	// Get existing member
	memberGet := &model.Member{SlackID: "1234"}
	err = dal.GetMemberBySlackID(memberGet)
	assert.Nil(t, err)
	assert.Equal(t, member.Name, memberGet.Name)

	// Update existing member
	memberUpdated := &model.Member{
		SlackID:  "1234",
		Name:     "Little Bruno",
		Position: "A REAL GUY",
	}
	err = dal.UpdateMember(memberUpdated)
	assert.Nil(t, err)
	err = dal.GetMemberBySlackID(memberGet)
	assert.Nil(t, err)
	assert.Equal(t, memberUpdated.Name, memberGet.Name)
	assert.Equal(t, memberUpdated.Position, memberGet.Position)

	// Delete existing member
	err = dal.DeleteMember(&model.Member{SlackID: "1234"})
	assert.Nil(t, err)
	err = dal.GetMemberBySlackID(&model.Member{SlackID: "1234"})
	assert.NotNil(t, err)
}
