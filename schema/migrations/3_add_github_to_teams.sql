# Destructive to team data
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;

CREATE TABLE teams (
    name TEXT UNIQUE,
    github_team_id INTEGER PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() at time zone 'utc')
);

CREATE TABLE team_members (
    team_github_team_id INTEGER REFERENCES teams(github_team_id) ON DELETE CASCADE,
    member_slack_id TEXT REFERENCES members(slack_id) ON DELETE CASCADE,
    PRIMARY KEY (team_github_team_id, member_slack_id)
);
