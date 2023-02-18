CREATE DATABASE greenlight;

\c greenlight

CREATE ROLE greenlight WITH LOGIN PASSWORD `pa55word`;

CREATE EXTENSION IF NOT EXISTS citext;