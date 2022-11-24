CREATE TABLE public.users
(
    id                  TEXT NOT NULL PRIMARY KEY,
    email               TEXT NOT NULL,
    firstname           TEXT NOT NULL,
    lastname            TEXT NOT NULL,
    display_name        TEXT NOT NULL,
    registration_date   TIMESTAMPTZ NOT NULL,
    hashed_password     BYTEA NOT NULL,
    roles               INT[] NOT NULL
);