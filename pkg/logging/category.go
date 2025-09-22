package logging

type Category string
type SubCategory string
type ExtraKey string

const (
	General  Category = "General"
	IO       Category = "IO"
	Internal Category = "Internal"

	Http            Category = "Http"
	JWT             Category = "JWT"
	Postgres        Category = "Postgres"
	Snowflake       Category = "SnowFlake"
	S3              Category = "S3"
	Redis           Category = "Redis"
	Validation      Category = "Validation"
	RequestResponse Category = "RequestResponse"
	Prometheus      Category = "Prometheus"
	Rabbit          Category = "Rabbit"
)

const (
	// General
	Startup         SubCategory = "Startup"
	Shutdown        SubCategory = "Shutdown"
	ExternalService SubCategory = "ExternalService"

	// Postgres
	Migration SubCategory = "Migration"
	Select    SubCategory = "Select"
	Rollback  SubCategory = "Rollback"
	CreateTx  SubCategory = "CreateTx"
	Update    SubCategory = "Update"
	Delete    SubCategory = "Delete"
	Insert    SubCategory = "Insert"
	Upsert    SubCategory = "Upsert"

	//S3 Storage
	IsExist       SubCategory = "IsExist"
	CreatedID     SubCategory = "CreateID"
	GenerateToken SubCategory = "GenerateToken"
	// Internal
	Api                 SubCategory = "Api"
	HashPassword        SubCategory = "HashPassword"
	DefaultRoleNotFound SubCategory = "DefaultRoleNotFound"
	FailedToCreateUser  SubCategory = "FailedToCreateUser"

	// Validation
	MobileValidation   SubCategory = "MobileValidation"
	PasswordValidation SubCategory = "PasswordValidation"

	// Http
	HttpError SubCategory = "HttpError"
	// IO
	RemoveFile SubCategory = "RemoveFile"

	//Rabbit
	Publish  SubCategory = "Publish"
	Received SubCategory = "Received"
)

const (
	AppName      ExtraKey = "AppName"
	LoggerName   ExtraKey = "Logger"
	ClientIp     ExtraKey = "ClientIp"
	HostIp       ExtraKey = "HostIp"
	Method       ExtraKey = "Method"
	StatusCode   ExtraKey = "StatusCode"
	BodySize     ExtraKey = "BodySize"
	BodyData     ExtraKey = "BodyData"
	Path         ExtraKey = "Path"
	Latency      ExtraKey = "Latency"
	RequestBody  ExtraKey = "RequestBody"
	ResponseBody ExtraKey = "ResponseBody"
	ErrorMessage ExtraKey = "ErrorMessage"
)
