package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/traPtitech/knoQ/utils"

	"github.com/traPtitech/knoQ/router"
	"github.com/traPtitech/knoQ/router/service"

	repo "github.com/traPtitech/knoQ/repository"

	"github.com/carlescere/scheduler"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

var (
	SESSION_KEY = []byte(os.Getenv("SESSION_KEY")[:32])
)

func main() {
	db, err := repo.SetupDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	googleAPI := &repo.GoogleAPIRepository{
		CalendarID: os.Getenv("TRAQ_CALENDARID"),
	}
	bytes, _ := ioutil.ReadFile("service.json")
	googleAPI.Config, err = google.JWTConfigFromJSON(bytes, calendar.CalendarReadonlyScope)
	if err == nil {
		googleAPI.Setup()
	}

	logger, _ := zap.NewDevelopment()
	handler := &router.Handlers{
		Dao: service.Dao{
			Repo: &repo.GormRepository{
				DB:       db,
				TokenKey: SESSION_KEY,
			},
			InitExternalUserGroupRepo: func(token string, ver repo.TraQVersion) interface {
				repo.UserBodyRepository
				repo.GroupRepository
			} {
				traQRepo := new(repo.TraQRepository)
				traQRepo.Token = token
				traQRepo.Version = ver
				traQRepo.Host = "https://q.trap.jp/api"
				traQRepo.NewRequest = traQRepo.DefaultNewRequest
				return traQRepo
			},
			InitTraPGroupRepo: func(token string, ver repo.TraQVersion) interface {
				repo.GroupRepository
			} {
				traPGroupRepo := new(repo.TraPGroupRepository)
				traPGroupRepo.Token = token
				traPGroupRepo.Version = ver
				traPGroupRepo.Host = "https://q.trap.jp/api"
				traPGroupRepo.NewRequest = traPGroupRepo.DefaultNewRequest
				return traPGroupRepo
			},
			ExternalRoomRepo: googleAPI,
		},
		Logger:     logger,
		SessionKey: SESSION_KEY,
		ClientID:   os.Getenv("CLIENT_ID"),
		SessionOption: sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 30,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
		WebhookID:         os.Getenv("WEBHOOK_ID"),
		WebhookSecret:     os.Getenv("WEBHOOK_SECRET"),
		ActivityChannelID: os.Getenv("CHANNEL_ID"),
		Origin:            os.Getenv("ORIGIN"),
	}

	e := handler.SetupRoute(db)

	// webhook
	job := utils.InitPostEventToTraQ(handler.Repo, handler.WebhookSecret,
		handler.ActivityChannelID, handler.WebhookID, handler.Origin)
	scheduler.Every().Day().At("08:00").Run(job)

	// サーバースタート
	go func() {
		if err := e.Start(":3000"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
