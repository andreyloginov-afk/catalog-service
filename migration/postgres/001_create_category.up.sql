CREATE TABLE category (
                          id          BIGSERIAL UNIQUE NOT NULL,
                          guid        UUID PRIMARY KEY,
                          name        VARCHAR(255) NOT NULL UNIQUE,
                          created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
                          updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);