CREATE TABLE product (
                         id            BIGSERIAL UNIQUE NOT NULL,
                         guid          UUID PRIMARY KEY,
                         name          VARCHAR(255) NOT NULL UNIQUE,
                         description   TEXT,
                         price         DECIMAL(12,2) NOT NULL,
                         category_guid UUID NOT NULL,
                         created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
                         updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

                         CONSTRAINT fk_product_category
                             FOREIGN KEY (category_guid)
                                 REFERENCES category (guid)
                                 ON DELETE RESTRICT
);