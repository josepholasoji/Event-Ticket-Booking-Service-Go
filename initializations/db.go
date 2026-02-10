package initializations

import (
	"fmt"
	"os"
	"sync"

	// Import mysql driver
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/kafka-go"
)

var (
	// mysqlDB type
	MySQLDB *sql.DB

	// KafkaProducer type
	KafkaProducerEvents *kafka.Writer
	KafkaConsumerEvents *kafka.Reader
	// KafkaProducer type
	KafkaProducerNotification *kafka.Writer
	KafkaConsumerNotification *kafka.Reader
)

func ConnectionToMySQLDB() error {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	fmt.Printf("\nConnecting to MySQL server with the following %s\n", dsn)

	retryLimit := 5
	retryCount := 0
	for {
		sqlDB, err := sql.Open("mysql", dsn)
		if err != nil {
			retryCount++
			if retryCount >= retryLimit {
				return fmt.Errorf("failed to connect to MySQL after %d attempts: %v", retryCount, err)
			}
			fmt.Printf("MySQL connection failed (attempt %d/%d): %v. Retrying...\n", retryCount, retryLimit, err)
			time.Sleep(time.Duration(1000*retryCount) * time.Millisecond)
			continue
		}
		MySQLDB = sqlDB
		break
	}

	// Connection pool settings (important!)
	MySQLDB.SetMaxOpenConns(25)
	MySQLDB.SetMaxIdleConns(25)
	MySQLDB.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := MySQLDB.Ping(); err != nil {
		return err
	}
	return nil
}

func CloseMySQLDB() error {
	if MySQLDB != nil {
		return MySQLDB.Close()
	}
	return nil
}

func ConnectToRedis() error {
	// Implement Redis connection logic here
	return nil
}

func CreateProducer(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(os.Getenv("KAFKA_BROKERS")),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func CreateConsumer(topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKERS")},
		Topic:   topic,
		GroupID: os.Getenv("KAFKA_GROUP_ID"),
	})
}

var once sync.Once

func Initialize() {
	once.Do(func() {
		ConnectionToMySQLDB()
		KafkaProducerEvents = CreateProducer("events")
		KafkaConsumerEvents = CreateConsumer("events")
	})
}
