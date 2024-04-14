package tests

import (
	"BannerService/internal/consts"
	"BannerService/internal/gateway"
	"BannerService/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestOptions struct {
	Banner          models.Banner
	TagId           int32
	FeatureId       int32
	UseLastRevision bool
	Token           string
	ExpectedStatus  int
	HasRequestBody  bool
}

func CreateGetBannerRequest(url, token string) (*http.Request, error) {
	// Создаем запрос GET на "/user_banner"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("token", token)
	return req, nil
}

func BannerAssertion(t *testing.T, banner1, banner2 models.Banner) {

	assert.Equal(t, banner1.Id, banner2.Id)
	assert.Equal(t, banner1.FeatureId, banner2.FeatureId)

	assert.True(t, banner1.TagIds.Equal(banner2.TagIds))
	assert.Equal(t, banner1.Content, banner2.Content)

	assert.Equal(t, banner1.CreatedAt.Format("2006-01-02 15:04:05"),
		banner2.CreatedAt.Format("2006-01-02 15:04:05"))
	assert.Equal(t, banner1.UpdatedAt.Format("2006-01-02 15:04:05"),
		banner2.UpdatedAt.Format("2006-01-02 15:04:05"))

	assert.Equal(t, *banner1.IsActive, *banner2.IsActive)
	assert.Equal(t, banner1.Version, banner2.Version)
}

func (s *APITestSuite) FindInCache(banner models.Banner) {
	ctx := context.Background()

	var temp models.Banner

	for _, tagId := range banner.TagIds {
		err := s.cache.Get(ctx, gateway.GenKey(tagId, banner.FeatureId)).Scan(&temp)
		if err != nil {
			s.Fail("Error with finding banner in cache. ",
				"id: %d, key: %s", banner.Id, gateway.GenKey(tagId, banner.FeatureId))
		}

		BannerAssertion(s.T(), banner, temp)
	}

}

func (s *APITestSuite) GetBannerTesting(options TestOptions) {

	tagId := options.TagId
	featureId := options.FeatureId
	UseLastRevision := options.UseLastRevision

	req, err := CreateGetBannerRequest(
		fmt.Sprintf("/user_banner?tag_id=%d&feature_id=%d&use_last_revision=%v",
			tagId, featureId, UseLastRevision),
		options.Token)

	if err != nil {
		s.Fail("Error with creating request")
	}

	rr := httptest.NewRecorder()
	s.server.Handler.ServeHTTP(rr, req)

	assert.Equal(s.T(), options.ExpectedStatus, rr.Code)

	if options.HasRequestBody {
		var actualResponse models.Banner
		if err := json.NewDecoder(rr.Body).Decode(&actualResponse); err != nil {
			s.Fail("Failed to decode response body", err)
		}

		BannerAssertion(s.T(), options.Banner, actualResponse)
	}

}

func (s *APITestSuite) TestGetUserBannerPositive() {

	banner := banners[0]

	options := []TestOptions{
		{
			Banner:          banner,
			TagId:           banner.TagIds[0],
			FeatureId:       banner.FeatureId,
			UseLastRevision: true,
			Token:           consts.UserToken,
			ExpectedStatus:  http.StatusOK,
			HasRequestBody:  true,
		},
		{
			Banner:          banner,
			TagId:           banner.TagIds[0],
			FeatureId:       banner.FeatureId,
			UseLastRevision: true,
			Token:           consts.AdminToken,
			ExpectedStatus:  http.StatusOK,
			HasRequestBody:  true,
		},
	}

	for _, option := range options {
		s.GetBannerTesting(option)
	}

}

func (s *APITestSuite) TestGetUserBannerWithoutValidToken() {

	options := []TestOptions{
		{
			Banner:          models.Banner{},
			TagId:           1,
			FeatureId:       1,
			UseLastRevision: false,
			Token:           "",
			ExpectedStatus:  http.StatusUnauthorized,
			HasRequestBody:  false,
		},
		{
			Banner:          models.Banner{},
			TagId:           1,
			FeatureId:       1,
			UseLastRevision: false,
			Token:           "fflfkfl",
			ExpectedStatus:  http.StatusForbidden,
			HasRequestBody:  false,
		},
		{
			Banner:          models.Banner{},
			TagId:           1,
			FeatureId:       1,
			UseLastRevision: false,
			Token:           "1",
			ExpectedStatus:  http.StatusForbidden,
			HasRequestBody:  false,
		},
	}

	for _, option := range options {
		s.GetBannerTesting(option)
	}

}

func (s *APITestSuite) TestGetUserBannerBadRequest() {

	options := []TestOptions{
		{
			Banner:          models.Banner{},
			TagId:           0,
			FeatureId:       1,
			UseLastRevision: false,
			Token:           consts.UserToken,
			ExpectedStatus:  http.StatusBadRequest,
			HasRequestBody:  false,
		},
		{
			Banner:          models.Banner{},
			TagId:           1,
			FeatureId:       0,
			UseLastRevision: false,
			Token:           consts.UserToken,
			ExpectedStatus:  http.StatusBadRequest,
			HasRequestBody:  false,
		},
		{
			Banner:          models.Banner{},
			TagId:           -1,
			FeatureId:       -1,
			UseLastRevision: false,
			Token:           consts.UserToken,
			ExpectedStatus:  http.StatusBadRequest,
			HasRequestBody:  false,
		},
	}

	for _, option := range options {
		s.GetBannerTesting(option)
	}

}

func (s *APITestSuite) TestGetUserBannerNotFound() {

	status := http.StatusNotFound

	options := []TestOptions{
		{
			Banner:          models.Banner{},
			TagId:           5,
			FeatureId:       5,
			UseLastRevision: true,
			Token:           consts.UserToken,
			ExpectedStatus:  status,
			HasRequestBody:  false,
		},
		{
			Banner:          models.Banner{},
			TagId:           10000,
			FeatureId:       1,
			UseLastRevision: true,
			Token:           consts.UserToken,
			ExpectedStatus:  status,
			HasRequestBody:  false,
		},
		{
			Banner:          models.Banner{},
			TagId:           1,
			FeatureId:       1000,
			UseLastRevision: true,
			Token:           consts.UserToken,
			ExpectedStatus:  status,
			HasRequestBody:  false,
		},
	}

	for _, option := range options {
		s.GetBannerTesting(option)
	}

}

func (s *APITestSuite) TestGetUserBannerNotActive() {

	banner := models.Banner{
		Id:        3,
		Version:   1,
		TagIds:    models.Tags{5, 6},
		FeatureId: 4,
		Content:   models.ModelMap{"ggkg": float64(5), "vv": "flfllf"},
		IsActive:  boolPtr(false),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.populateDB(banner)
	if err != nil {
		s.FailNow("Failed to populate DB", err)
	}

	options := []TestOptions{
		{
			Banner:          models.Banner{},
			TagId:           banner.TagIds[0],
			FeatureId:       banner.FeatureId,
			UseLastRevision: true,
			Token:           consts.UserToken,
			ExpectedStatus:  http.StatusNotFound,
			HasRequestBody:  false,
		},
		{
			Banner:          banner,
			TagId:           banner.TagIds[0],
			FeatureId:       banner.FeatureId,
			UseLastRevision: true,
			Token:           consts.AdminToken,
			ExpectedStatus:  http.StatusOK,
			HasRequestBody:  true,
		},
	}

	for _, option := range options {
		s.GetBannerTesting(option)
	}

}

func (s *APITestSuite) TestGetUserBannerNewVersion() {

	banner := banners[0]

	s.GetBannerTesting(TestOptions{
		Banner:          banner,
		TagId:           banner.TagIds[0],
		FeatureId:       banner.FeatureId,
		UseLastRevision: true,
		Token:           consts.UserToken,
		ExpectedStatus:  http.StatusOK,
		HasRequestBody:  true,
	})

	newBanner := banner
	newBanner.Version = 2
	newBanner.UpdatedAt = time.Now()
	err := s.populateDB(newBanner)
	if err != nil {
		s.FailNow("Failed to populate DB", err)
	}

	_, err = s.db.Exec("UPDATE banners SET is_active=$1 WHERE id=$2 AND version=$3",
		false, banner.Id, banner.Version)

	if err != nil {
		s.FailNow("Failed to update DB", err)
	}

	options := []TestOptions{
		{
			Banner:          banner,
			TagId:           newBanner.TagIds[0],
			FeatureId:       newBanner.FeatureId,
			UseLastRevision: false,
			Token:           consts.UserToken,
			ExpectedStatus:  http.StatusOK,
			HasRequestBody:  true,
		},
		{
			Banner:          banner,
			TagId:           newBanner.TagIds[0],
			FeatureId:       newBanner.FeatureId,
			UseLastRevision: false,
			Token:           consts.AdminToken,
			ExpectedStatus:  http.StatusOK,
			HasRequestBody:  true,
		},
		// Тут берем актуальную версию
		{
			Banner:          newBanner,
			TagId:           newBanner.TagIds[0],
			FeatureId:       newBanner.FeatureId,
			UseLastRevision: true,
			Token:           consts.UserToken,
			ExpectedStatus:  http.StatusOK,
			HasRequestBody:  true,
		},
		{
			Banner:          newBanner,
			TagId:           newBanner.TagIds[0],
			FeatureId:       newBanner.FeatureId,
			UseLastRevision: true,
			Token:           consts.AdminToken,
			ExpectedStatus:  http.StatusOK,
			HasRequestBody:  true,
		},
	}

	for _, option := range options {
		s.GetBannerTesting(option)
	}

}

func (s *APITestSuite) TestGetUserBannerNewVersionWithOtherData() {

	banner := banners[0]

	s.GetBannerTesting(TestOptions{
		Banner:          banner,
		TagId:           banner.TagIds[0],
		FeatureId:       banner.FeatureId,
		UseLastRevision: true,
		Token:           consts.UserToken,
		ExpectedStatus:  http.StatusOK,
		HasRequestBody:  true,
	})

	newBanner := banner
	newBanner.Version = 2
	newBanner.TagIds = models.Tags{20, 25}
	newBanner.UpdatedAt = time.Now()

	err := s.populateDB(newBanner)
	if err != nil {
		s.FailNow("Failed to populate DB", err)
	}

	_, err = s.db.Exec("UPDATE banners SET is_active=$1 WHERE id=$2 AND version=$3",
		false, banner.Id, banner.Version)

	if err != nil {
		s.FailNow("Failed to update DB", err)
	}

	options := []TestOptions{
		{ // Пробуем получить старый баннер(доступен ещё 5 минут)
			Banner:          banner,
			TagId:           banner.TagIds[0],
			FeatureId:       banner.FeatureId,
			UseLastRevision: false,
			Token:           consts.UserToken,
			ExpectedStatus:  http.StatusOK,
			HasRequestBody:  true,
		},
		{
			Banner:          banner,
			TagId:           banner.TagIds[0],
			FeatureId:       banner.FeatureId,
			UseLastRevision: false,
			Token:           consts.AdminToken,
			ExpectedStatus:  http.StatusOK,
			HasRequestBody:  true,
		},
		// Тут берем актуальную версию

		{ // Пробуем получить старый баннер напрямую из бд
			Banner:          models.Banner{},
			TagId:           banner.TagIds[0],
			FeatureId:       banner.FeatureId,
			UseLastRevision: true,
			Token:           consts.UserToken,
			ExpectedStatus:  http.StatusNotFound,
			HasRequestBody:  false,
		},
		{
			Banner:          models.Banner{},
			TagId:           banner.TagIds[0],
			FeatureId:       banner.FeatureId,
			UseLastRevision: true,
			Token:           consts.AdminToken,
			ExpectedStatus:  http.StatusNotFound,
			HasRequestBody:  false,
		},

		{ // Пробуем получить новый баннер
			Banner:          newBanner,
			TagId:           newBanner.TagIds[0],
			FeatureId:       newBanner.FeatureId,
			UseLastRevision: false,
			Token:           consts.UserToken,
			ExpectedStatus:  http.StatusOK,
			HasRequestBody:  false,
		},
		{
			Banner:          newBanner,
			TagId:           newBanner.TagIds[0],
			FeatureId:       newBanner.FeatureId,
			UseLastRevision: false,
			Token:           consts.AdminToken,
			ExpectedStatus:  http.StatusOK,
			HasRequestBody:  false,
		},
	}

	for _, option := range options {
		s.GetBannerTesting(option)
	}
}

func (s *APITestSuite) TestGetUserBannerNewVersionWithActive1Version() {

	banner := banners[0]

	//s.GetBannerTesting(TestOptions{
	//	Banner:          banner,
	//	TagId:           banner.TagIds[0],
	//	FeatureId:       banner.FeatureId,
	//	UseLastRevision: true,
	//	Token:           consts.UserToken,
	//	ExpectedStatus:  http.StatusOK,
	//	HasRequestBody:  true,
	//})

	newBanner := banner
	newBanner.Version = 2
	newBanner.UpdatedAt = time.Now()
	newBanner.IsActive = boolPtr(false)

	err := s.populateDB(newBanner)
	if err != nil {
		s.FailNow("Failed to populate DB", err)
	}

	options := []TestOptions{
		{
			Banner:          banner,
			TagId:           newBanner.TagIds[0],
			FeatureId:       newBanner.FeatureId,
			UseLastRevision: true,
			Token:           consts.UserToken,
			ExpectedStatus:  http.StatusOK,
			HasRequestBody:  true,
		},
		{
			Banner:          banner,
			TagId:           newBanner.TagIds[0],
			FeatureId:       newBanner.FeatureId,
			UseLastRevision: true,
			Token:           consts.AdminToken,
			ExpectedStatus:  http.StatusOK,
			HasRequestBody:  true,
		},
	}

	for _, option := range options {
		s.GetBannerTesting(option)
	}

}

func (s *APITestSuite) TestGetUserBannerDeactivatedAllVersion() {

	banner := banners[0]

	s.GetBannerTesting(TestOptions{
		Banner:          banner,
		TagId:           banner.TagIds[0],
		FeatureId:       banner.FeatureId,
		UseLastRevision: true,
		Token:           consts.UserToken,
		ExpectedStatus:  http.StatusOK,
		HasRequestBody:  true,
	})

	newBanner := banner
	newBanner.Version = 2
	newBanner.UpdatedAt = time.Now()
	newBanner.IsActive = boolPtr(false)

	err := s.populateDB(newBanner)
	if err != nil {
		s.FailNow("Failed to populate DB", err)
	}

	_, err = s.db.Exec("UPDATE banners SET is_active=$1 WHERE id=$2",
		false, banner.Id)

	options := []TestOptions{
		{
			Banner:          banner,
			TagId:           newBanner.TagIds[0],
			FeatureId:       newBanner.FeatureId,
			UseLastRevision: false,
			Token:           consts.UserToken,
			ExpectedStatus:  http.StatusOK,
			HasRequestBody:  true,
		},
		{
			Banner:          banner,
			TagId:           newBanner.TagIds[0],
			FeatureId:       newBanner.FeatureId,
			UseLastRevision: false,
			Token:           consts.AdminToken,
			ExpectedStatus:  http.StatusOK,
			HasRequestBody:  true,
		},
		{
			Banner:          models.Banner{},
			TagId:           newBanner.TagIds[0],
			FeatureId:       newBanner.FeatureId,
			UseLastRevision: true,
			Token:           consts.UserToken,
			ExpectedStatus:  http.StatusNotFound,
			HasRequestBody:  false,
		},
		{
			Banner:          newBanner,
			TagId:           newBanner.TagIds[0],
			FeatureId:       newBanner.FeatureId,
			UseLastRevision: true,
			Token:           consts.AdminToken,
			ExpectedStatus:  http.StatusOK,
			HasRequestBody:  true,
		},
	}

	for _, option := range options {
		s.GetBannerTesting(option)
	}

}
