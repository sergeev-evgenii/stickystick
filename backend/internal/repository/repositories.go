package repository

import "gorm.io/gorm"

type Repositories struct {
	User     UserRepository
	Video    VideoRepository
	Comment  CommentRepository
	Like     LikeRepository
	Category CategoryRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:     NewUserRepository(db),
		Video:    NewVideoRepository(db),
		Comment:  NewCommentRepository(db),
		Like:     NewLikeRepository(db),
		Category: NewCategoryRepository(db),
	}
}
