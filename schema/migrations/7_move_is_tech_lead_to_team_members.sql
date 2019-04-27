ALTER TABLE members
DROP is_tech_lead;

ALTER TABLE team_members
ADD is_tech_lead BOOLEAN DEFAULT false;
