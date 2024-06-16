package database

import (
	"carthero/internal/model"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// GetRiders returns a list of riders from the database.
	GetRiders() ([]model.Rider, error)

	// fetch free riders

	GetFreeRiders() ([]model.Rider, error)

	// update rider status
	UpdateRiderStatus(id int, status bool) error

	// create rider
	CreateRider(rider model.Rider) (model.Rider, error)

	// delete rider
	DeleteRider(id int) error
}

type service struct {
	db *gorm.DB
}

var (
	dbname     = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return &service{
			db: dbInstance.db,
		}
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	dbInstance = &service{db: db}
	log.Printf("Connected to database: %s", dbname)
	db.AutoMigrate(new(model.Rider))

	return dbInstance

}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {

	stats := make(map[string]string)

	db, err := s.db.DB()
	if err != nil {
		stats["status"] = "error"
		stats["error"] = err.Error()
	}
	err = db.Ping()
	if err != nil {
		stats["status"] = "error"
		stats["error"] = err.Error()
	}
	stats["status"] = "ok"
	stats["database"] = dbname

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dbname)
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

// GetRiders returns a list of riders from the database.
func (s *service) GetRiders() ([]model.Rider, error) {
	var riders []model.Rider
	result := s.db.Find(&riders)
	return riders, result.Error
}

// GetFreeRiders returns a list of free riders from the database.
func (s *service) GetFreeRiders() ([]model.Rider, error) {
	var riders []model.Rider
	result := s.db.Where("assigned = ?", false).Find(&riders)
	return riders, result.Error
}

// UpdateRiderStatus updates the status of a rider in the database.
func (s *service) UpdateRiderStatus(id int, status bool) error {
	result := s.db.Model(&model.Rider{}).Where("id = ?", id).Update("assigned", status)
	return result.Error
}

// CreateRider creates a new rider in the database.
func (s *service) CreateRider(rider model.Rider) (model.Rider, error) {
	result := s.db.Create(&rider)

	return rider, result.Error
}

// DeleteRider deletes a rider from the database.
func (s *service) DeleteRider(id int) error {
	result := s.db.Where("id = ?", id).Delete(&model.Rider{})
	return result.Error
}
