package env

import (
	"github.com/spf13/viper"
)

type Config struct {
	GrpcPort               string
	HttpPort               string
	PasswordResetEndPoint  string
	AccountConfirmEndPoint string
	StudentSignupEndpoint  string
	SignupLinkLife         int
	ResetPasswordLinkLife  int
	JWTSignupTokenSecret   string
	StripeFreePlan         string
	Database               *DatabaseConfig
	AuthService            *AuthServiceConfig
	Noreply                *STMPCredentiels
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

type AuthServiceConfig struct {
	Host string
	Port string
}

type STMPCredentiels struct {
	Host     string
	Port     int
	Username string
	Password string
	Email    string
	APIKey   string
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

	viper.BindEnv("STMP_HOST")
	// viper.SetDefault("STMP_PORT", 465)
	viper.BindEnv("STMP_PORT")
	viper.BindEnv("NO_REPLY_EMAIL")
	viper.BindEnv("SMTP_USERNAME")
	viper.BindEnv("SMPT_PASSWORD")
	viper.BindEnv("SEND_GRID_API_KEY")

	viper.BindEnv("RESET_PASSWORD_ENDPOINT")
	viper.BindEnv("ACCOUNT_CONFIRM_ENDPOINT")
	viper.BindEnv("AGENT_SIGNUP_ENDPOINT")
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
		AccountConfirmEndPoint: viper.GetString("ACCOUNT_CONFIRM_ENDPOINT"),
		StudentSignupEndpoint:  viper.GetString("AGENT_SIGNUP_ENDPOINT"),
		SignupLinkLife:         viper.GetInt("SIGNUP_LINK_LIFE"),
		JWTSignupTokenSecret:   viper.GetString("JWT_SIGNUP_TOKEN_SECRET"),
		StripeFreePlan:         viper.GetString("STRIPE_FREE_PLAN"),

		AuthService: &AuthServiceConfig{
			Host: viper.GetString("AUTH_SERVICE_HOST"),
			Port: viper.GetString("AUTH_SERVICE_PORT"),
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
		Noreply: &STMPCredentiels{
			Host:     viper.GetString("STMP_HOST"),
			Port:     viper.GetInt("STMP_PORT"),
			Username: viper.GetString("NO_REPLY_USERNAME"),
			Password: viper.GetString("NO_REPLY_PASSWORD"),
			Email:    viper.GetString("NO_REPLY_EMAIL"),
			APIKey:   viper.GetString("SEND_GRID_API_KEY"),
		},
	}
}
