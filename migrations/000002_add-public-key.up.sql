ALTER TABLE public.tokens
ADD COLUMN public_key TEXT NOT NULL UNIQUE;

DROP INDEX IF EXISTS tokens_refresh_index;

CREATE UNIQUE INDEX tokens_public_key_index ON public.tokens (public_key);