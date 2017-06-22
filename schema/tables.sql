CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS members;
CREATE TABLE members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(128),
    first_name VARCHAR(128),
    last_name VARCHAR(128),
    github_username VARCHAR(128),
    program VARCHAR(64),
    image_url VARCHAR(256),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() at time zone 'utc')
);

DROP TABLE IF EXISTS teams;
CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(128),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() at time zone 'utc')
);

DROP TABLE IF EXISTS team_members;
CREATE TABLE team_members (
    team_id UUID REFERENCES teams(id),
    member_id UUID REFERENCES members(id),
    PRIMARY KEY (team_id, member_id)
);
