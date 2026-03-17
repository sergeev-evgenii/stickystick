-- Таблица жалоб на видео (если не создаётся через RUN_MIGRATIONS=1)
CREATE TABLE IF NOT EXISTS video_reports (
    id BIGSERIAL PRIMARY KEY,
    video_id BIGINT NOT NULL,
    user_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_video_reports_video_id ON video_reports(video_id);
CREATE INDEX IF NOT EXISTS idx_video_reports_user_id ON video_reports(user_id);
CREATE INDEX IF NOT EXISTS idx_video_reports_deleted_at ON video_reports(deleted_at);
