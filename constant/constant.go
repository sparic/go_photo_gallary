package constant

const (
	// JWT constants
	JWT_SECRET          = "JWT_SECRET"
	JWT                 = "jwt"
	JWT_EXP_MINUTE      = 30
	PHOTO_STORAGE_ADMIN = "admin"

	// Server constants
	SERVER_PORT = "6666"
	PAGE_SIZE   = 20

	// DB constants
	// DB_CONNECT = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"
	DB_CONNECT = "%s:%s@%s:%s/%s?charset=utf8&parseTime=True&loc=Local"
	DB_TYPE    = "mysql"
	DB_HOST    = "localhost"
	DB_PORT    = "3306"
	DB_USER    = "root"
	DB_PWD     = "root"
	DB_NAME    = "root"

	// Auth constants
	COOKIE_MAX_AGE = 1800
	LOGIN_MAX_AGE  = 1800
	LOGIN_USER     = "LOGIN_"
	// Redis constants
	REDIS_HOST = "localhost"
	REDIS_PORT = "6379"
	// COS constants
	// COS_BUCKET_NAME = "COS_BUCKET_NAME"
	// COS_APP_ID      = "COS_APP_ID"
	// COS_REGION      = "COS_REGION"
	// COS_SECRET_ID   = "COS_SECRET_ID"
	// COS_SECRET_KEY  = "COS_SECRET_KEY"

	SERVER_DOMAIN = "localhost"
	SERVER_PATH   = "/"
)
