package model

// Root represent a config at root level.
type Root struct {
	Server         Server         `mapstructure:"server"`
	Service        Service        `mapstructure:"service"`
	DefaultAccount DefaultAccount `mapstructure:"defaultAccount"`
	Token          Token          `mapstructure:"token"`
}

// Server provide info for http server
type Server struct {
	Name            string `mapstructure:"name"`
	Version         string `mapstructure:"version"`
	Port            int    `mapstructure:"port"`
	ShutdownTimeout int    `mapstructure:"shutdownTimeout"`
	Heroku          bool   `mapstructure:"heroku"`
	Logging         bool   `mapstructure:"logging"`
	Monitoring      bool   `mapstructure:"monitoring"`
	HTTPS           bool   `mapstructure:"https"`
}

// Service struct provide all services
type Service struct {
	Database Database `mapstructure:"database"`
	Email    Email    `mapstructure:"email"`
	NewRelic NewRelic `mapstructure:"newrelic"`
	Fluent   Fluent   `mapstructure:"fluent"`
	Kafka    Kafka    `mapstructure:"kafka"`
	RabbitMQ RabbitMQ `mapstructure:"rabbitmq"`
	PGO      PGO      `mapstructure:"pgo"`
	MQTT     MQTT     `mapstructure:"mqtt"`
}

// DefaultAccount provide default account info
type DefaultAccount struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// Token stuct provide configuration for a jwt authentication
type Token struct {
	Accesstoken  Accesstoken  `mapstructure:"accessToken"`
	Refreshtoken Refreshtoken `mapstructure:"refreshToken"`
	OTP          OTP          `mapstructure:"otp"`
}

// Accesstoken provide configuration for a jwt authentication
type Accesstoken struct {
	Enable     bool   `mapstructure:"enable"`
	PublicKey  string `mapstructure:"publicKey"`
	PrivateKey string `mapstructure:"privateKey"`
	JWTTimeout int    `mapstructure:"jwtTimeout"`
}

// Refreshtoken provide configuration for a refresh token
type Refreshtoken struct {
	Enable     bool   `mapstructure:"enable"`
	PublicKey  string `mapstructure:"publicKey"`
	PrivateKey string `mapstructure:"privateKey"`
	JWTTimeout int    `mapstructure:"jwtTimeout"`
}

// OTP provide configuration for a OTP token
type OTP struct {
	SecretKey string `mapstructure:"secretKey"`
	Timeout   int    `mapstructure:"timeout"`
}

// Database provide info for database connection.
type Database struct {
	OneDrive   OneDrive   `mapstructure:"onedrive"`
	Dropbox    Dropbox    `mapstructure:"dropbox"`
	Drive      Drive      `mapstructure:"drive"`
	Sharepoint Sharepoint `mapstructure:"sharepoint"`
	MongoDB    MongoDB    `mapstructure:"mongodb"`
	Postgres   Postgres   `mapstructure:"postgres"`
	Redis      Redis      `mapstructure:"redis"`
	Collection Collection `mapstructure:"collection"`
}

// OneDrive provide a connection information for onedrive
type OneDrive struct {
	URL          string `mapstructure:"url"`
	AccessToken  string `mapstructure:"accessToken"`
	RefreshToken string `mapstructure:"refreshToken"`
}

// Dropbox provide a connection information for dropbox
type Dropbox struct {
	Token string `mapstructure:"token"`
}

// Drive provide a connection information for drive
type Drive struct {
	APIKey string `mapstructure:"apiKey"`
}

// Sharepoint provide a connection information for sharepoint
type Sharepoint struct {
	SiteURL  string `mapstructure:"siteURL"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// MongoDB provide a connection information for mongodb
type MongoDB struct {
	User     string   `mapstructure:"user"`
	Password string   `mapstructure:"password"`
	Hosts    []string `mapstructure:"hosts"`
	DB       string   `mapstructure:"db"`
	Options  []string `mapstructure:"options"`
}

// Postgres provide a connection information for postgres
type Postgres struct {
	Host     []string `mapstructure:"host"`
	User     string   `mapstructure:"user"`
	Password string   `mapstructure:"password"`
	DBName   string   `mapstructure:"dbName"`
	Port     string   `mapstructure:"port"`
}

// Redis stuct provide info for redis
type Redis struct {
	Prefix   string `mapstructure:"prefix"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	DB       int    `mapstructure:"db"`
}

// Collection provide database's collection for mongodb
type Collection struct {
	Document string `mapstructure:"document"`
	User     string `mapstructure:"user"`
}

// Email provide config for email service
type Email struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Key      string `mapstructure:"key"`
	From     string `mapstructure:"from"`
	Subject  string `mapstructure:"subject"`
	Message  string `mapstructure:"message"`
}

// NewRelic provide config for New Relic service
type NewRelic struct {
	LicenseKey string `mapstructure:"licenseKey"`
}

// PGO provide config for New Relic service
type PGO struct {
	SpaceName string `mapstructure:"spaceName"`
}

// Fluent provide config for New Relic service
type Fluent struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	Prefix string `mapstructure:"prefix"`
	Tag    string `mapstructure:"tag"`
}

// Kafka provide config for New Relic service
type Kafka struct {
	Hosts     []string `mapstructure:"hosts"`
	Partition int32    `mapstructure:"partition"`
}

// RabbitMQ provide config for New Relic service
type RabbitMQ struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Route    string `mapstructure:"route"`
}

// MQTT provide config for New Relic service
type MQTT struct {
	Host      string `mapstructure:"host"`
	QoS       byte   `mapstructure:"qos"`
	PrefixPub string `mapstructure:"prefixPub"`
	PrefixSub string `mapstructure:"prefixSub"`
}
