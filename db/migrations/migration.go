package migrations

import (
	"fmt"

	"github.com/CRS-Project/crs-backend/internal/entity"
	mylog "github.com/CRS-Project/crs-backend/internal/pkg/logger"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	fmt.Println(mylog.ColorizeInfo("\n=========== Start Migrate ==========="))
	mylog.Infof("Migrating Tables...")

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return err
	}

	//migrate table
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Package{},
		&entity.UserDiscipline{},
		&entity.Comment{},
		&entity.DisciplineGroup{},
		&entity.DisciplineGroupConsolidator{},
		&entity.DisciplineListDocument{},
		&entity.DisciplineListDocumentConsolidator{},
	); err != nil {
		return err
	}

	return nil
}
