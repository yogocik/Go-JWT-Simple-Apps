package utility

import "gorm.io/gorm"

func HandleGormError(db *gorm.DB) error {
	if db.Error != nil {
		return db.Error
	}
	return nil
}
