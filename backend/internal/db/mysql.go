package db

import (
	"time"

	"github.com/rs/zerolog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"interview-coach/backend/internal/config"
	"interview-coach/backend/internal/model"
)

// Open connects to MySQL with GORM and returns a configured database handle.
//
// Inputs: cfg contains the DB_DSN and pool settings; logger is used for
// structured GORM output.
// Outputs: an initialized *gorm.DB, or an error when the connection fails.
func Open(cfg config.Config, logger zerolog.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DBDSN), &gorm.Config{
		Logger: gormlogger.New(zerologWriter{logger: logger}, gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		}),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	return db, nil
}

// AutoMigrate creates or updates the database schema for the backend models.
//
// Inputs: db is an active GORM database connection.
// Outputs: nil when migration succeeds, otherwise the migration error.
func AutoMigrate(db *gorm.DB) error {
	return db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(&model.InterviewSession{})
}

type zerologWriter struct {
	logger zerolog.Logger
}

// Printf writes GORM log messages through zerolog so database logs remain structured.
//
// Inputs: format is the GORM log template; args are substituted into the format.
// Outputs: none. The formatted message is emitted as a structured log line.
func (w zerologWriter) Printf(format string, args ...interface{}) {
	w.logger.Warn().Msgf(format, args...)
}
