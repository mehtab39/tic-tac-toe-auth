


CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);


GRANT ALL PRIVILEGES ON TABLE users TO mehtabgill;

GRANT USAGE, SELECT ON SEQUENCE users_id_seq TO mehtabgill;
