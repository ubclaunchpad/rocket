package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubclaunchpad/rocket/model"
)

func TestMemberCreateGetUpdateRemove(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	dal, cleanupFunc, err := newTestDBConnection()
	assert.Nil(t, err)
	defer cleanupFunc()

	// Create a new member
	member := &model.Member{
		SlackID: "1234",
		Name:    "Big Bruno",
	}
	err = dal.CreateMember(member)
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

func TestSetMemberIsAdmin(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	dal, cleanupFunc, err := newTestDBConnection()
	assert.Nil(t, err)
	defer cleanupFunc()

	// Create a new member
	member := &model.Member{
		SlackID: "1234",
		Name:    "Big Bruno",
		IsAdmin: true,
	}
	err = dal.CreateMember(member)
	assert.Nil(t, err)

	// Set member admin status
	err = dal.SetMemberIsAdmin(&model.Member{SlackID: "1234"})
	assert.Nil(t, err)

	// Get member
	memberGet := &model.Member{SlackID: "1234"}
	err = dal.GetMemberBySlackID(memberGet)
	assert.Nil(t, err)
	assert.False(t, memberGet.IsAdmin)
}
