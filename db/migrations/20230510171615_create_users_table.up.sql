CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    uuid VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) UNIQUE,
    email_verified_at BIGINT,
    created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
    updated_at BIGINT
);

-- create or replace function public.handle_new_user()
-- returns trigger as $$
-- declare
-- username varchar;
-- begin
--   username := split_part(new.email, '@', 1);
-- insert into public.users (uuid, email, username)
-- values (new.id, new.email, username);
-- return new;
-- end;
-- $$ language plpgsql security definer;

-- create trigger on_auth_user_created
--     after insert on auth.users
--     for each row execute procedure public.handle_new_user();
