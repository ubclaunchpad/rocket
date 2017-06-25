DROP TABLE IF EXISTS members;
CREATE TABLE members (
    slack_id TEXT PRIMARY KEY,
    name TEXT,
    email TEXT UNIQUE,
    github_username TEXT,
    program TEXT,
    position TEXT,
    image_url TEXT,
    is_admin BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() at time zone 'utc')
);

DROP TABLE IF EXISTS teams;
CREATE TABLE teams (
    name TEXT PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() at time zone 'utc')
);

DROP TABLE IF EXISTS team_members;
CREATE TABLE team_members (
    team_name TEXT REFERENCES teams(name) ON DELETE CASCADE,
    member_email TEXT REFERENCES members(email) ON DELETE CASCADE,
    PRIMARY KEY (team_name, member_email)
);
