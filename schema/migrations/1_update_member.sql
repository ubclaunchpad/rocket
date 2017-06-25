ALTER TABLE members
ADD position TEXT;

ALTER TABLE members
ADD is_admin BOOLEAN DEFAULT false; 