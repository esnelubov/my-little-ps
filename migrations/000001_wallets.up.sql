CREATE TABLE IF NOT EXISTS wallets
(
    id bigserial,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    name text COLLATE pg_catalog."default" NOT NULL,
    country text COLLATE pg_catalog."default" NOT NULL,
    city text COLLATE pg_catalog."default" NOT NULL,
    currency text COLLATE pg_catalog."default" NOT NULL,
    balance bigint NOT NULL DEFAULT 0,
    worker integer NOT NULL DEFAULT 0,
    CONSTRAINT wallets_pkey PRIMARY KEY (id)
);

CREATE INDEX mlp_wallets_worker ON wallets (worker);