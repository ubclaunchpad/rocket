package data

import "github.com/ubclaunchpad/rocket/model"

// GetMemberBySlackID populates the given member with information from the DB
// or returns an error.
func (dal *DAL) GetMemberBySlackID(member *model.Member) error {
	return dal.db.Model(member).
		Where("slack_id = ?slack_id").
		Select()
}

// GetMembers populates the given members with information for all members from
// the DB or returns an error.
func (dal *DAL) GetMembers(members *model.Members) error {
	return dal.db.Model(members).
		Order("name ASC").
		Select()
}

// GetTechLeads populates given members with all current tech leads
func (dal *DAL) GetTechLeads(members *model.Members) error {
	return dal.db.Model(members).
		Where("is_tech_lead = 't'").
		Order("name ASC").
		Select()
}

// GetAdmins populates the given members with information for all admin members
// or returns an error.
func (dal *DAL) GetAdmins(members *model.Members) error {
	return dal.db.Model(members).
		Where("is_admin = 't'").
		Order("name ASC").
		Select()
}

// CreateMember adds the given member to the DB or returns an error.
func (dal *DAL) CreateMember(member *model.Member) error {
	_, err := dal.db.Model(member).
		OnConflict("DO NOTHING").
		Insert()
	return err
}

// UpdateMember updates the entry in the DB for the given member or returns an
// error. The name and email fields will be updated as long as a member exists
// with `member`'s SlackID. The position field will only be updated if the
// existing row has an empty value for that column. No other columns will be
// altered.
func (dal *DAL) UpdateMember(member *model.Member) error {
	existing := &model.Member{SlackID: member.SlackID}
	if err := dal.GetMemberBySlackID(existing); err != nil {
		return err
	}
	// As long as we have member name and email, update them
	if member.Name != "" {
		existing.Name = member.Name
	}
	if member.Email != "" {
		existing.Email = member.Email
	}
	// Only update position if we don't already have one
	if existing.Position == "" {
		existing.Position = member.Position
	}
	_, err := dal.db.Model(existing).
		Update(
			"name", existing.Name,
			"email", existing.Email,
			"position", existing.Position)
	return err
}

// DeleteMember deletes a member from the DB or returns an error.
func (dal *DAL) DeleteMember(member *model.Member) error {
	_, err := dal.db.Model(member).
		Where("slack_id = ?slack_id").
		Delete()
	return err
}

// SetMemberName updates the name of the given member in the DB or returns
// an error.
func (dal *DAL) SetMemberName(member *model.Member) error {
	_, err := dal.db.Model(member).
		Set("name = ?name").
		Update()

	return err
}

// SetMemberEmail updates the name of the given member in the DB or returns
// an error.
func (dal *DAL) SetMemberEmail(member *model.Member) error {
	_, err := dal.db.Model(member).
		Set("email = ?email").
		Update()

	return err
}

// SetMemberGitHubUsername updates the GitHub username of the given member in
// the DB or returns an error.
func (dal *DAL) SetMemberGitHubUsername(member *model.Member) error {
	_, err := dal.db.Model(member).
		Set("github_username = ?github_username").
		Update()

	return err
}

// SetMemberMajor updates the major of the given member in the DB or returns
// an error.
func (dal *DAL) SetMemberMajor(member *model.Member) error {
	_, err := dal.db.Model(member).
		Set("program = ?program").
		Update()

	return err
}

// SetMemberPosition updates the position of the given member in the DB or returns
// an error.
func (dal *DAL) SetMemberPosition(member *model.Member) error {
	_, err := dal.db.Model(member).
		Set("position = ?position").
		Update()

	return err
}

// SetMemberBiography updates the bio of the given member in the DB or returns
// an error.
func (dal *DAL) SetMemberBiography(member *model.Member) error {
	_, err := dal.db.Model(member).
		Set("biography = ?biography").
		Update()

	return err
}

// SetMemberImageURL updates the image URL of the given member in the DB or returns
// an error.
func (dal *DAL) SetMemberImageURL(member *model.Member) error {
	_, err := dal.db.Model(member).
		Set("image_url = ?image_url").
		Update()

	return err
}

// SetMemberIsAdmin updates whether the given member is an admin in the DB
// or returns an error.
func (dal *DAL) SetMemberIsAdmin(member *model.Member) error {
	_, err := dal.db.Model(member).
		Set("is_admin = ?is_admin").
		Update()

	return err
}

// SetMemberIsTechLead updates whether the given member is a tech lead in
// the DB or returns an error.
func (dal *DAL) SetMemberIsTechLead(member *model.Member) error {
	_, err := dal.db.Model(member).
		Set("is_tech_lead = ?is_tech_lead").
		Update()

	return err
}
