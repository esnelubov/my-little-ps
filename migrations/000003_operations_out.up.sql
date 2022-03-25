CREATE TABLE IF NOT EXISTS out_operations
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
    CONSTRAINT out_operations_pkey PRIMARY KEY (id)
);

CREATE INDEX mlp_out_operations_transaction_id ON out_operations (transaction_id);

CREATE INDEX mlp_out_operations_origin_wallet_id ON out_operations (origin_wallet_id);

CREATE INDEX mlp_out_operations_target_wallet_id ON out_operations (target_wallet_id);

CREATE INDEX mlp_out_operations_status ON out_operations (status);