DROP INDEX IF EXISTS tokens_public_key_index;

ALTER TABLE public.tokens
DROP COLUMN IF EXISTS public_key;

-- Recreate the old unique index for the refresh_token column
CREATE UNIQUE INDEX tokens_refresh_index ON public.tokens (refresh_token);