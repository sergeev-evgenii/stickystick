package repository

import (
	"log"
	"sticky-stick/backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	// Проверяем, существует ли таблица videos
	var tableExists bool
	err := db.Raw(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.tables 
			WHERE table_name = 'videos'
		)
	`).Scan(&tableExists).Error
	
	if err == nil && tableExists {
		// Проверяем, существует ли старая колонка video_url
		var hasVideoURL bool
		err = db.Raw(`
			SELECT EXISTS (
				SELECT 1 
				FROM information_schema.columns 
				WHERE table_name = 'videos' 
				AND column_name = 'video_url'
			)
		`).Scan(&hasVideoURL).Error
		
		if err == nil && hasVideoURL {
			log.Println("Migrating: Renaming video_url to media_url...")
			// Переименовываем video_url в media_url
			if err := db.Exec("ALTER TABLE videos RENAME COLUMN video_url TO media_url").Error; err != nil {
				log.Printf("Warning: Failed to rename video_url column: %v", err)
			} else {
				log.Println("Successfully renamed video_url to media_url")
			}
		}
		
		// Проверяем, существует ли колонка media_type
		var hasMediaType bool
		err = db.Raw(`
			SELECT EXISTS (
				SELECT 1 
				FROM information_schema.columns 
				WHERE table_name = 'videos' 
				AND column_name = 'media_type'
			)
		`).Scan(&hasMediaType).Error
		
		if err == nil && !hasMediaType {
			log.Println("Migrating: Adding media_type column...")
			if err := db.Exec("ALTER TABLE videos ADD COLUMN media_type VARCHAR(20) DEFAULT 'video'").Error; err != nil {
				log.Printf("Warning: Failed to add media_type column: %v", err)
			} else {
				log.Println("Successfully added media_type column")
			}
		}
		
		// Обновляем существующие записи, у которых media_type NULL
		db.Exec("UPDATE videos SET media_type = 'video' WHERE media_type IS NULL")
		
		// Проверяем, существует ли колонка moderation_status
		var hasModerationStatus bool
		err = db.Raw(`
			SELECT EXISTS (
				SELECT 1 
				FROM information_schema.columns 
				WHERE table_name = 'videos' 
				AND column_name = 'moderation_status'
			)
		`).Scan(&hasModerationStatus).Error
		
		if err == nil && !hasModerationStatus {
			log.Println("Migrating: Adding moderation_status column...")
			if err := db.Exec("ALTER TABLE videos ADD COLUMN moderation_status VARCHAR(20) DEFAULT 'pending'").Error; err != nil {
				log.Printf("Warning: Failed to add moderation_status column: %v", err)
			} else {
				log.Println("Successfully added moderation_status column")
				// Устанавливаем статус 'approved' для существующих видео
				db.Exec("UPDATE videos SET moderation_status = 'approved' WHERE moderation_status IS NULL")
			}
		}
		
		// Проверяем, существует ли колонка is_admin в таблице users
		var hasIsAdmin bool
		err = db.Raw(`
			SELECT EXISTS (
				SELECT 1 
				FROM information_schema.columns 
				WHERE table_name = 'users' 
				AND column_name = 'is_admin'
			)
		`).Scan(&hasIsAdmin).Error
		
		if err == nil && !hasIsAdmin {
			log.Println("Migrating: Adding is_admin column to users...")
			if err := db.Exec("ALTER TABLE users ADD COLUMN is_admin BOOLEAN DEFAULT FALSE").Error; err != nil {
				log.Printf("Warning: Failed to add is_admin column: %v", err)
			} else {
				log.Println("Successfully added is_admin column to users")
			}
		}

		// Нормализация media_url и thumbnail_url к относительному формату (videos/xxx.mp4)
		log.Println("Migrating: Normalizing media_url to relative paths...")
		db.Exec(`UPDATE videos SET media_url = substring(media_url from 10) WHERE media_url LIKE '/uploads/%'`)
		db.Exec(`
			UPDATE videos SET media_url = substring(
				regexp_replace(media_url, '^https?://[^/]+', '') from 10
			) WHERE media_url LIKE 'http%' AND media_url LIKE '%/uploads/%'
		`)
		db.Exec(`UPDATE videos SET thumbnail_url = substring(thumbnail_url from 10) WHERE thumbnail_url LIKE '/uploads/%'`)
		db.Exec(`
			UPDATE videos SET thumbnail_url = substring(
				regexp_replace(thumbnail_url, '^https?://[^/]+', '') from 10
			) WHERE thumbnail_url LIKE 'http%' AND thumbnail_url LIKE '%/uploads/%'
		`)
	}

	// Запускаем AutoMigrate для всех моделей
	log.Println("Running AutoMigrate...")
	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Video{},
		&models.Comment{},
		&models.Like{},
		&models.ActivityLog{},
		&models.SiteSettings{},
	); err != nil {
		return err
	}

	// Создаём запись настроек по умолчанию, если её нет
	var count int64
	if err := db.Model(&models.SiteSettings{}).Count(&count).Error; err == nil && count == 0 {
		if err := db.Create(&models.SiteSettings{
			ID:                    1,
			ShowViewCount:         true,
			// VK-посты: ссылка на Telegram + ссылка на проект
			DefaultPublishVK:       "мы в телеграмм\nhttps://t.me/uncensored_mems\nhttps://stickystick.ru",
			// Telegram-посты: ссылка на VK + ссылка на проект
			DefaultPublishTelegram: "мы в вк\nhttps://vk.com/club236352692\nhttps://stickystick.ru",
			DefaultPublishMax:      "мы в макс ссылка на макс\nhttps://stickystick.ru",
		}).Error; err != nil {
			log.Printf("Warning: Failed to create default site_settings: %v", err)
		} else {
			log.Println("Created default site_settings (show_view_count=true)")
		}
	}

	// Если запись уже есть (обновление), проставляем дефолты для новых полей, если они пустые
	db.Exec(`UPDATE site_settings SET default_publish_vk = 'мы в телеграмм
https://t.me/uncensored_mems
https://stickystick.ru' WHERE (default_publish_vk IS NULL OR default_publish_vk = '') AND id = 1`)
	db.Exec(`UPDATE site_settings SET default_publish_telegram = 'мы в вк
https://vk.com/club236352692
https://stickystick.ru' WHERE (default_publish_telegram IS NULL OR default_publish_telegram = '') AND id = 1`)
	db.Exec(`UPDATE site_settings SET default_publish_max = 'мы в макс ссылка на макс
https://stickystick.ru' WHERE (default_publish_max IS NULL OR default_publish_max = '') AND id = 1`)

	// Создаем категорию "юмор", если её нет
	var humorCategory models.Category
	result := db.Where("slug = ?", "humor").First(&humorCategory)
	if result.Error != nil {
		humorCategory = models.Category{
			Name: "Юмор",
			Slug: "humor",
		}
		if err := db.Create(&humorCategory).Error; err != nil {
			log.Printf("Warning: Failed to create humor category: %v", err)
		} else {
			log.Println("Created humor category")
		}
	}
	
	log.Println("Migrations completed successfully")
	return nil
}
