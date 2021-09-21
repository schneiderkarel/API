DROP INDEX user_id_index;
ALTER TABLE user_service.users DROP CONSTRAINT user_id_pk;

DROP TABLE user_service.users;
DROP SCHEMA IF EXISTS user_service;
