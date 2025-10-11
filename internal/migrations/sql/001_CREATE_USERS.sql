DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'user') THEN
       CREATE DATABASE users;
END IF;
END $$;