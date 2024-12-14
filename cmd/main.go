package main

import (
	"os"
	"time"

	"github.com/emersion/go-smtp"
	r "github.com/end1essrage/listener-smtp-rmq/rmq"
	s "github.com/end1essrage/listener-smtp-rmq/smtp"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

var (
	Env string

	SmtpAddress string
	SmtpDomain  string

	AmqpUrl          string
	AmqpQueueName    string
	AmqpExchangeName string
)

const (
	ENV_DEBUG = "ENV_DEBUG" //Для локального запуска в дебаг режиме
	ENV_LOCAL = "ENV_LOCAL" //Для локального запуска
	ENV_POD   = "ENV_POD"   //Для запуска в контейнере

	SMTP_ADDRESS = "SMTP_ADDRESS"
	SMTP_DOMAIN  = "SMTP_DOMAIN"

	AMQP_URL      = "AMQP_URL"
	AMQP_QUEUE    = "AMQP_QUEUE"
	AMQP_EXCHANGE = "AMQP_EXCHANGE"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	Env = os.Getenv("ENV")
	if Env == "" {
		if err := godotenv.Load(); err != nil {
			logrus.Warning("error while reading environment %s", err.Error())
		}
	}

	Env = os.Getenv("ENV")
	if Env == "" {
		logrus.Warn("cant set environment, setting to local by default")
		Env = ENV_LOCAL
	}

	logrus.Info("ENVIRONMENT IS " + Env)
	loadEnv()
}

func loadEnv() {
	SmtpAddress = os.Getenv(SMTP_ADDRESS)
	logrus.Info("SmtpAddress = " + SmtpAddress)
	SmtpDomain = os.Getenv(SMTP_DOMAIN)
	logrus.Info("SmtpDomain = " + SmtpDomain)

	AmqpUrl = os.Getenv(AMQP_URL)
	AmqpQueueName = os.Getenv(AMQP_QUEUE)
	AmqpExchangeName = os.Getenv(AMQP_EXCHANGE)
}

func main() {
	// конфигурация шифрования и аутентификации
	cfg := amqp.Config{}
	cfg.Vhost = "/"

	rClient := r.NewClient(AmqpUrl, AmqpQueueName, AmqpExchangeName, cfg)
	be := s.NewServer(rClient)

	s := smtp.NewServer(be)

	s.Addr = SmtpAddress
	s.Domain = SmtpDomain
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	logrus.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		logrus.Fatal(err)
	}
}
