SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

CREATE SCHEMA IF NOT EXISTS "user_service";

CREATE TABLE "user_service"."users" (
    "user_id" UUID NOT NULL,
    "name" VARCHAR(100) NOT NULL,
    "age" SMALLINT NOT NULL
);

ALTER TABLE "user_service"."users" ADD CONSTRAINT "user_id_pk" PRIMARY KEY ("user_id");
CREATE UNIQUE INDEX "user_id_index" ON "user_service"."users" USING btree ("user_id");
