CREATE TABLE IF NOT EXISTS operations
(
    id bigserial,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone,
    uid text COLLATE pg_catalog."default" NOT NULL,
    wallet_id bigint NOT NULL,
    type text COLLATE pg_catalog."default" NOT NULL,
    amount bigint NOT NULL,
    status text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT operations_pkey PRIMARY KEY (id)
)
