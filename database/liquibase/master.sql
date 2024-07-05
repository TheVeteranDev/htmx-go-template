--liquibase formatted sql

--changeset master:1
CREATE TABLE IF NOT EXISTS hellos (  
  id INT  PRIMARY KEY,
  word VARCHAR(255) NOT NULL,
  language VARCHAR(255) NOT NULL
);

--changeset master:2
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL
);