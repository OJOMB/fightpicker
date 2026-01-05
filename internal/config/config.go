package config

type Config struct {
	Env      string `mapstructure:"env"`
	Domain   string `mapstructure:"domain"`
	Port     int    `mapstructure:"port"`
	LogLevel int    `mapstructure:"log_level"`

	API           APIConfig           `mapstructure:"api"`
	Database      DatabaseConfig      `mapstructure:"database"`
	Cache         CacheConfig         `mapstructure:"cache"`
	Auth          AuthConfig          `mapstructure:"auth"`
	Observability ObservabilityConfig `mapstructure:"observability"`
	AWS           AWSConfig           `mapstructure:"aws"`
	EventBroker   EventBrokerConfig   `mapstructure:"event_broker"`
	Email         EmailConfig         `mapstructure:"email"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type CacheConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type APIConfig struct {
	Users    bool `mapstructure:"users"`
	Fighters bool `mapstructure:"fighters"`
}

type AuthConfig struct {
	PrivateKey      string `mapstructure:"private_key"`
	PublicKey       string `mapstructure:"public_key"`
	HashingCost     int    `mapstructure:"hashing_cost"`
	AccessTTLHours  int    `mapstructure:"access_ttl_hours"`
	RefreshTTLHours int    `mapstructure:"refresh_ttl_hours"`
	TokenAudience   string `mapstructure:"token_audience"`
	TokenIssuer     string `mapstructure:"token_issuer"`
	SSLMode         string `mapstructure:"ssl_mode"`
}

type ObservabilityConfig struct {
	OTel      OTelConfig      `mapstructure:"otel"`
	Pyroscope PyroscopeConfig `mapstructure:"pyroscope"`
}

type OTelConfig struct {
	Enable   bool   `mapstructure:"enable"`
	Endpoint string `mapstructure:"endpoint"`
}

type PyroscopeConfig struct {
	Enable   bool   `mapstructure:"enable"`
	Endpoint string `mapstructure:"endpoint"`
}

type AWSConfig struct {
	Region             string `mapstructure:"region"`
	AccessKeyID        string `mapstructure:"access_key_id"`
	SecretAccessKey    string `mapstructure:"secret_access_key"`
	S3Endpoint         string `mapstructure:"s3_endpoint"`
	S3MediaBucket      string `mapstructure:"s3_media_bucket"`
	S3MediaBucketRaw   string `mapstructure:"s3_media_bucket_raw"`
	UsePathStyle       bool   `mapstructure:"use_path_style"`
	PresignedPutURLTTL int    `mapstructure:"presigned_put_url_ttl_minutes"`
	PresignedGetURLTTL int    `mapstructure:"presigned_get_url_ttl_minutes"`
}

type EventBrokerConfig struct {
	SeedBrokers               []string `mapstructure:"seed_brokers"`
	GroupID                   string   `mapstructure:"group_id"`
	TopicProfilePictureUpload string   `mapstructure:"topic_profile_picture_upload"`
	TopicPostUserCreated      string   `mapstructure:"topic_post_user_created"`
	TopicPostUserVerified     string   `mapstructure:"topic_post_user_verified"`
}

type EmailConfig struct {
	SMTPHost       string `mapstructure:"smtp_host"`
	SMTPPort       int    `mapstructure:"smtp_port"`
	SMTPUser       string `mapstructure:"smtp_user"`
	SMTPPassword   string `mapstructure:"smtp_password"`
	AddressNoReply string `mapstructure:"address_no_reply"`
	SkipTLS        bool   `mapstructure:"skip_tls"`
}
