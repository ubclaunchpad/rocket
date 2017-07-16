ALTER TABLE team_members
DROP COLUMN member_email;

ALTER TABLE team_members
ADD COLUMN member_slack_id TEXT REFERENCES members(slack_id) ON DELETE CASCADE;

ALTER TABLE team_members
ADD PRIMARY KEY (team_name, member_slack_id);