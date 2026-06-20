package db

import (
	"fmt"
	"log"
	"time"

	"event_ticket_booking/config"
	refreshTokenEntity "event_ticket_booking/infrastructure/db/refresh_token/entity"
	refreshTokenRepo "event_ticket_booking/infrastructure/db/refresh_token/repository"
	userEntity "event_ticket_booking/infrastructure/db/user/entity"
	userRepo "event_ticket_booking/infrastructure/db/user/repository"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Db struct {
	UserRepo         userRepo.IRepository
	RefreshTokenRepo refreshTokenRepo.IRepository
}

func NewDb(config config.Config) Db {
	db := NewDbConnection(config.Db)

	// Auto-migrate the auth tables so the service is runnable out of the box.
	if err := db.AutoMigrate(&userEntity.Entity{}, &refreshTokenEntity.Entity{}); err != nil {
		panic(fmt.Sprintf("failed to auto-migrate database: %v", err))
	}

	return Db{
		UserRepo:         userRepo.NewRepository(db),
		RefreshTokenRepo: refreshTokenRepo.NewRepository(db),
	}
}

func NewDbConnection(config config.DbConfig) *gorm.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true&timeout=10s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DbName,
	)

	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("failed to get database instance: %v", err))
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	log.Printf("Connect Database successfully")
	return db
}
