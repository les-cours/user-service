package env

import (
	"github.com/spf13/viper"
)

type Config struct {
	GrpcPort               string
	HttpPort               string
	PasswordResetEndPoint  string
	TeacherConfirmEndPoint string
	StudentSignupEndpoint  string
	SignupLinkLife         int
	ResetPasswordLinkLife  int
	JWTSignupTokenSecret   string
	Database               *DatabaseConfig
	AuthService            *ServiceConfig
	LearningService        *ServiceConfig
	NoreplyEmail           string
}

type DatabaseConfig struct {
	PSQLConfig
}

type PSQLConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
	SslMode  string
}

type ServiceConfig struct {
	Host string
	Port string
}

var Settings *Config

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	viper.BindEnv("GRPC_PORT")
	viper.BindEnv("HTTP_PORT")

	//	viper.SetDefault("MONGO_DATABASE", "test")
	viper.BindEnv("MONGO_VISITORS_COLLECTION")
	viper.BindEnv("MONGO_ROOMS_COLLECTION")

	viper.BindEnv("POSTGRES_HOST")
	viper.BindEnv("POSTGRES_PORT")
	viper.BindEnv("POSTGRES_USERNAME")
	viper.BindEnv("POSTGRES_PASSWORD")
	viper.BindEnv("POSTGRES_DBNAME")
	viper.BindEnv("POSTGRES_SSL_MODE")

	viper.BindEnv("MONGO_HOST")
	viper.BindEnv("MONGO_REPLICASET")
	viper.BindEnv("MONGO_USERNAME")
	viper.BindEnv("MONGO_PASSWORD")
	viper.BindEnv("MONGO_DBNAME")
	viper.BindEnv("STAGE")

	viper.BindEnv("AUTH_SERVICE_HOST")
	viper.BindEnv("AUTH_SERVICE_PORT")

	viper.BindEnv("SETTINGS_SERVICE_HOST")
	viper.BindEnv("SETTINGS_SERVICE_PORT")

	viper.BindEnv("STRIPE_SERVICE_HOST")
	viper.BindEnv("STRIPE_SERVICE_PORT")
	viper.BindEnv("STRIPE_FREE_PLAN")

	viper.BindEnv("NO_REPLY_EMAIL")

	viper.BindEnv("RESET_PASSWORD_ENDPOINT")
	viper.BindEnv("APP_TEACHER_CONFIRM_ENDPOINT")
	viper.BindEnv("USER_SIGNUP_ENDPOINT")
	viper.BindEnv("SIGNUP_LINK_LIFE")         // 30day
	viper.BindEnv("RESET_PASSWORD_LINK_LIFE") //1800s == 30min
	viper.BindEnv("JWT_SIGNUP_TOKEN_SECRET")

	viper.BindEnv("VISITOR_SERVICE_HOST")
	viper.BindEnv("VISITOR_SERVICE_PORT")

	Settings = &Config{
		GrpcPort:               viper.GetString("GRPC_PORT"),
		HttpPort:               viper.GetString("HTTP_PORT"),
		PasswordResetEndPoint:  viper.GetString("RESET_PASSWORD_ENDPOINT"),
		ResetPasswordLinkLife:  viper.GetInt("RESET_PASSWORD_LINK_LIFE"),
		TeacherConfirmEndPoint: viper.GetString("TEACHER_CONFIRM_ENDPOINT"),
		StudentSignupEndpoint:  viper.GetString("USER_SIGNUP_ENDPOINT"),
		SignupLinkLife:         viper.GetInt("SIGNUP_LINK_LIFE"),
		JWTSignupTokenSecret:   viper.GetString("JWT_SIGNUP_TOKEN_SECRET"),
		NoreplyEmail:           viper.GetString("NO_REPLY_EMAIL"),

		AuthService: &ServiceConfig{
			Host: viper.GetString("AUTH_SERVICE_HOST"),
			Port: viper.GetString("AUTH_SERVICE_PORT"),
		},

		LearningService: &ServiceConfig{
			Host: viper.GetString("LEARNING_SERVICE_HOST"),
			Port: viper.GetString("LEARNING_SERVICE_PORT"),
		},

		Database: &DatabaseConfig{
			PSQLConfig{
				Host:     viper.GetString("POSTGRES_HOST"),
				Port:     viper.GetInt("POSTGRES_PORT"),
				Username: viper.GetString("POSTGRES_USERNAME"),
				Password: viper.GetString("POSTGRES_PASSWORD"),
				DbName:   viper.GetString("POSTGRES_DBNAME"),
				SslMode:  viper.GetString("POSTGRES_SSL_MODE"),
			},
		},
	}
}
