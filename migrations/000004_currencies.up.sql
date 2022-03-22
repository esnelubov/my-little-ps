CREATE TABLE IF NOT EXISTS currencies
(
    id bigserial,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone,
    name text COLLATE pg_catalog."default" NOT NULL UNIQUE,
    usd_rate bigint NOT NULL,
    CONSTRAINT currencies_pkey PRIMARY KEY (id)
)
