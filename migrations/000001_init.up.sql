CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS public.users
(
  guid        UUID NOT NULL,
  email       VARCHAR(255) NOT NULL,
  ip          INET NOT NULL,
  PRIMARY KEY (guid)
);


CREATE TABLE IF NOT EXISTS public.tokens
(
  guid        UUID REFERENCES users(guid) ON DELETE CASCADE NOT NULL,
  refresh_token       TEXT NOT NULL UNIQUE,
  PRIMARY KEY (guid)
);

CREATE UNIQUE INDEX tokens_refresh_index ON public.tokens (refresh_token);