-- create new permissions table that stores the permission levels
CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    code text NOT NULL
);

-- create new table users_permissions which denotes what permissions each user has
CREATE TABLE IF NOT EXISTS users_permissions (
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE ,
    permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);

-- add new permissions to the table
INSERT INTO permissions (code)
VALUES ('movies:read'), ('movies:write');