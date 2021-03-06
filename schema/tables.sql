DROP TABLE IF EXISTS members CASCADE;
CREATE TABLE members (
    slack_id TEXT PRIMARY KEY,
    name TEXT,
    email TEXT UNIQUE,
    github_username TEXT,
    program TEXT,
    position TEXT,
    biography TEXT,
    image_url TEXT,
    is_tech_lead BOOLEAN DEFAULT false,
    is_admin BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() at time zone 'utc')
);

DROP TABLE IF EXISTS teams CASCADE;
CREATE TABLE teams (
    name TEXT UNIQUE,
    github_team_id INTEGER PRIMARY KEY,
    platform TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() at time zone 'utc')
);

DROP TABLE IF EXISTS team_members CASCADE;
CREATE TABLE team_members (
    team_github_team_id INTEGER REFERENCES teams(github_team_id) ON DELETE CASCADE,
    member_slack_id TEXT REFERENCES members(slack_id) ON DELETE CASCADE,
    PRIMARY KEY (team_github_team_id, member_slack_id)
);
