package tests

import (
	"BannerService/internal/gateway"
	"BannerService/internal/models"
	"BannerService/internal/server"
	"BannerService/internal/service"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

var dbURI, redisURI string

func init() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file: ", err.Error())
	//}

	dbURI = os.Getenv("TEST_DB_URI")
	redisURI = os.Getenv("TEST_REDIS_URI")
}

type APITestSuite struct {
	suite.Suite

	db       *sqlx.DB
	cache    *redis.Client
	server   *http.Server
	services *service.Services
	gateways *gateway.Gateways
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupTest() {

	_, err := s.db.Exec("DELETE FROM banners")
	assert.NoError(s.T(), err)

	err = s.populateDB(banners...)
	assert.NoError(s.T(), err)

	err = s.cache.FlushAll(context.Background()).Err()
	assert.NoError(s.T(), err)
}

func (s *APITestSuite) SetupSuite() {
	db, err := gateway.NewPostgresDB(dbURI)
	if err != nil {
		println(dbURI)
		s.FailNow("Failed to connect to postgres", err)
	}
	s.db = db

	println(redisURI)
	cache := redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := cache.Ping(context.Background()).Err(); err != nil {
		s.FailNow("Failed to connect to redis", err)
	}
	s.cache = cache

	s.initDeps()
}

func (s *APITestSuite) initDeps() {
	gateways := gateway.NewGateway(s.db, s.cache, true)
	services := service.NewService(gateways)
	port := os.Getenv("TEST_PORT")

	s.server = server.NewHttpServer(":"+port, services)
	s.services = services
	s.gateways = gateways

	go s.HandleRequests()
}

func (s *APITestSuite) HandleRequests() {

	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			s.FailNow("Failed to start server ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

}

func (s *APITestSuite) populateDB(data ...models.Banner) error {
	for _, banner := range data {

		createQuery := "INSERT INTO banners (id, version, feature_id, content, is_active, created_at, update_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"

		_, err := s.db.Exec(createQuery, banner.Id, banner.Version, banner.FeatureId,
			banner.Content, banner.IsActive, banner.CreatedAt, banner.UpdatedAt)
		if err != nil {
			return err
		}

		createTagsRelQuery := "INSERT INTO tags_banners (banner_id, banner_version, tag_id) VALUES ($1, $2, $3)"
		for _, tagId := range banner.TagIds {
			if _, err := s.db.Exec(createTagsRelQuery, banner.Id, banner.Version, tagId); err != nil {
				return err
			}
		}

	}

	return nil
}

func (s *APITestSuite) TearDownSuite() {
	_ = s.db.Close()
	_ = s.cache.Close()

	p, _ := os.FindProcess(os.Getpid())
	_ = p.Signal(syscall.SIGINT)
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}
