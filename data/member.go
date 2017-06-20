package data

import "github.com/ubclaunchpad/rocket/model"

func (dal *DAL) GetMemberByID(member *model.Member) error {
	return dal.db.Model(member).
		Where("id = ?id").
		Select()
}

func (dal *DAL) GetMembers(members *model.Members) error {
	return dal.db.Model(members).Select()
}
