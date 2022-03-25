CREATE TABLE IF NOT EXISTS in_operations
(
    id bigserial,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone,
    operation_id text COLLATE pg_catalog."default" NOT NULL UNIQUE,
    transaction_id text COLLATE pg_catalog."default" NOT NULL,
    origin_wallet_id bigint,
    target_wallet_id bigint,
    amount bigint NOT NULL,
    currency text NOT NULL,
    status text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT in_operations_pkey PRIMARY KEY (id)
);

CREATE INDEX mlp_in_operations_transaction_id ON in_operations (transaction_id);

CREATE INDEX mlp_in_operations_origin_wallet_id ON in_operations (origin_wallet_id);

CREATE INDEX mlp_in_operations_target_wallet_id ON in_operations (target_wallet_id);

CREATE INDEX mlp_in_operations_status ON in_operations (status);