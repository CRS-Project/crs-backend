package seeders

import (
	"fmt"

	"github.com/CRS-Project/crs-backend/db/seeder/seeds"
	mylog "github.com/CRS-Project/crs-backend/internal/pkg/logger"
	"gorm.io/gorm"
)

func Seeding(db *gorm.DB) error {
	seeders := []func(*gorm.DB) error{
		seeds.SeederUser,
	}

	fmt.Println(mylog.ColorizeInfo("\n=========== Start Seeding ==========="))
	for _, seeder := range seeders {
		if err := seeder(db); err != nil {
			return err
		}
	}

	return nil
}
