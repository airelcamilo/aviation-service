CREATE TABLE IF NOT EXISTS airport (
	id SERIAL PRIMARY KEY,
    type VARCHAR(20),
    facility_name VARCHAR(100),
    faa VARCHAR(10) UNIQUE,
    icao VARCHAR(10) UNIQUE,
    region VARCHAR(50),
    state VARCHAR(50),
    county VARCHAR(50),
    city VARCHAR(50),
    ownership CHAR(2),
    use CHAR(2),
    manager VARCHAR(100),
    manager_phone VARCHAR(20),
    latitude VARCHAR(20),
    longitude VARCHAR(20),
    status VARCHAR(20)
);

