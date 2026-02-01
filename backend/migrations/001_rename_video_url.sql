-- Миграция для переименования video_url в media_url и добавления media_type
-- Выполните этот скрипт вручную, если автоматическая миграция не сработала

-- Переименовываем video_url в media_url (если колонка существует)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'videos' 
        AND column_name = 'video_url'
    ) THEN
        ALTER TABLE videos RENAME COLUMN video_url TO media_url;
        RAISE NOTICE 'Column video_url renamed to media_url';
    ELSE
        RAISE NOTICE 'Column video_url does not exist, skipping rename';
    END IF;
END $$;

-- Добавляем колонку media_type, если её нет
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'videos' 
        AND column_name = 'media_type'
    ) THEN
        ALTER TABLE videos ADD COLUMN media_type VARCHAR(20) DEFAULT 'video';
        RAISE NOTICE 'Column media_type added';
    ELSE
        RAISE NOTICE 'Column media_type already exists, skipping';
    END IF;
END $$;

-- Обновляем существующие записи
UPDATE videos SET media_type = 'video' WHERE media_type IS NULL;
