-- Колонка is_hidden в videos: скрыть видео из раздачи без удаления
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'videos' AND column_name = 'is_hidden'
    ) THEN
        ALTER TABLE videos ADD COLUMN is_hidden BOOLEAN NOT NULL DEFAULT FALSE;
        CREATE INDEX IF NOT EXISTS idx_videos_is_hidden ON videos(is_hidden);
        RAISE NOTICE 'Column is_hidden added to videos';
    END IF;
END $$;
