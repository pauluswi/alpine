-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS pinstore (
    "id" UUID NOT NULL PRIMARY KEY,
    "customer_id" VARCHAR NOT NULL,
    "status" VARCHAR NOT NULL DEFAULT 1,
    "credential" VARCHAR NOT NULL,
    "metadata" JSONB NOT NULL DEFAULT '{}',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- Add unique index 
CREATE UNIQUE INDEX IF NOT EXISTS idx_unq_customer_id ON pinstore (customer_id);
