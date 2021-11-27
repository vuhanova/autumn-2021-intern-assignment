CREATE TABLE IF NOT EXISTS users
(
    ID      BIGSERIAL PRIMARY KEY,
    balance real
);


CREATE TABLE IF NOT EXISTS transaction
(
    ID BIGSERIAL PRIMARY KEY,
    to_id BIGINT,
    from_id BIGINT,
    money REAL NOT NULL,
    created TIMESTAMP NOT NULL,
    FOREIGN KEY (to_id) REFERENCES users(ID),
    FOREIGN KEY (from_id) REFERENCES users(ID)
    );