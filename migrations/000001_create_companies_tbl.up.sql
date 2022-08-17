CREATE TABLE IF NOT EXISTS companies (
    name VARCHAR(100),
    code VARCHAR(100),
    country VARCHAR(100),
    website VARCHAR(100),
    phone VARCHAR(32),
    PRIMARY KEY (name, code)
);
