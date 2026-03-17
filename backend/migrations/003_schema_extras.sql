-- Колонки moderation_status (videos) и is_admin (users)
-- Выполнять только если таблицы уже есть, а колонок нет.

-- moderation_status в videos
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'videos' AND column_name = 'moderation_status'
    ) THEN
        ALTER TABLE videos ADD COLUMN moderation_status VARCHAR(20) DEFAULT 'pending';
        UPDATE videos SET moderation_status = 'approved' WHERE moderation_status IS NULL;
        RAISE NOTICE 'Column moderation_status added to videos';
    END IF;
END $$;

-- is_admin в users
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'users' AND column_name = 'is_admin'
    ) THEN
        ALTER TABLE users ADD COLUMN is_admin BOOLEAN DEFAULT FALSE;
        RAISE NOTICE 'Column is_admin added to users';
    END IF;
END $$;
