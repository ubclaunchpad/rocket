package data

import "github.com/ubclaunchpad/rocket/model"

func (dal *DAL) GetMemberBySlackID(member *model.Member) error {
	return dal.db.Model(member).
		Where("slack_id = ?slack_id").
		Select()
}

func (dal *DAL) GetMembers(members *model.Members) error {
	return dal.db.Model(members).Select()
}

func (dal *DAL) CreateMember(member *model.Member) error {
	_, err := dal.db.Model(member).
		OnConflict("DO NOTHING").
		Insert()
	return err
}

func (dal *DAL) SetMemberName(member *model.Member) error {
	_, err := dal.db.Model(member).
		Set("name = ?name").
		Update()

	return err
}

func (dal *DAL) SetMemberEmail(member *model.Member) error {
	_, err := dal.db.Model(member).
		Set("email = ?email").
		Update()

	return err
}

func (dal *DAL) SetMemberGitHubUsername(member *model.Member) error {
	_, err := dal.db.Model(member).
		Set("github_username = ?github_username").
		Update()

	return err
}

func (dal *DAL) SetMemberMajor(member *model.Member) error {
	_, err := dal.db.Model(member).
		Set("program = ?program").
		Update()

	return err
}

func (dal *DAL) SetMemberPosition(member *model.Member) error {
	_, err := dal.db.Model(member).
		Set("position = ?position").
		Update()

	return err
}

func (dal *DAL) SetMemberImageURL(member *model.Member) error {
	_, err := dal.db.Model(member).
		Set("image_url = ?image_url").
		Update()

	return err
}

func (dal *DAL) SetMemberIsAdmin(member *model.Member) error {
	_, err := dal.db.Model(member).
		Set("is_admin = ?is_admin").
		Update()

	return err
}
