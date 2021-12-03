CREATE TABLE statuses (
    id   BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO statuses (id, name) VALUES (1, 'PENDING'), (2, 'CREATED');