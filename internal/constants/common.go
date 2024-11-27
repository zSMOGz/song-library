package constants

const (
	DefaultFormat        = "%s%s"
	DefaultAddressFormat = "%s://%s:%s"

	ErrFormatAddition = ": %w"
	ErrFormat         = "%s: %w"

	SQLExtension = ".sql"
	SQLSuffix    = "_%s" + SQLExtension
	//БД
	PostgresConnectionString = "postgres://%s:%s@%s:%s/%s?sslmode=%s"
	PostgresDriver           = "postgres"

	// API Роутеры
	APIBasePath   = "/api"
	APISongsPath  = APIBasePath + "/songs"
	APIVersesPath = APIBasePath + "/verses"
	MetricsPath   = "/metrics"

	// Пути API для песен
	APISongDelete = APISongsPath + "/delete"
	APISongUpdate = APISongsPath + "/update"
	APISongCreate = APISongsPath + "/create"
	APISongInfo   = APISongsPath + "/info"

	// Пути
	ProjectRootPath = "../.."

	// Пути SQL
	VerseQueriesPath = "queries/verses"
	SongQueriesPath  = "queries/songs"
	// SQL Запросы на получение данных
	QueryGet              = "get"
	QueryCreateSong       = "create"
	QueryCreateSimpleSong = "create_simple"
	QueryUpdateSong       = "update"
	QueryDeleteSong       = "delete"
	QueryListSongs        = "list"

	// Поля логов
	LogFieldMethod   = "method"
	LogFieldPath     = "path"
	LogFieldStatus   = "status"
	LogFieldDuration = "duration"
	LogMsgRequest    = "Request processed"

	// Метрики
	MetricHTTPRequestsTotal = "http_requests_total"
	MetricHTTPRequestsHelp  = "Общее количество HTTP запросов"

	// Надписи для метрик
	MetricLabelMethod   = "method"
	MetricLabelEndpoint = "endpoint"
	MetricLabelStatus   = "status"

	// Параметры URL запроса
	QueryParamSongID   = "song_id"
	QueryParamPage     = "page"
	QueryParamPageSize = "page_size"

	// Значения по умолчанию
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 50

	// Заголовки
	HeaderContentType     = "Content-Type"
	HeaderContentTypeJSON = "application/json"
	HeaderCacheControl    = "Cache-Control"
	CacheControlValue     = "public, max-age=300"

	// Названия методов для обработчиков
	HandlerGetVerses   = "GetVerses"
	HandlerGetSongs    = "GetSongs"
	HandlerDeleteSong  = "DeleteSong"
	HandlerUpdateSong  = "UpdateSong"
	HandlerCreateSong  = "CreateSong"
	HandlerGetSongInfo = "GetSongInfo"

	// Параметры URL запроса
	QueryParamID      = "id"
	QueryParamTitle   = "title"
	QueryParamArtist  = "artist"
	QueryParamGroup   = "group"
	QueryParamSong    = "song"
	QueryParamAlbum   = "album"
	QueryParamGenre   = "genre"
	QueryParamYear    = "year"
	QueryParamPerPage = "per_page"

	// Формат URL API
	APIInfoURLFormat = APISongInfo + "?group=%s&song=%s"

	// Environment variables
	EnvDBHost         = "DB_HOST"
	EnvDBPort         = "DB_PORT"
	EnvDBUser         = "DB_USER"
	EnvDBPassword     = "DB_PASSWORD"
	EnvDBName         = "DB_NAME"
	EnvDBSSLMode      = "DB_SSLMODE"
	EnvServerHost     = "SERVER_HOST"
	EnvServerPort     = "SERVER_PORT"
	EnvServerProtocol = "SERVER_PROTOCOL"
	// Configuration files
	EnvFileName = ".env"

	// Swagger пути
	SwaggerPath    = "/swagger/"
	SwaggerDocPath = "/swagger/doc.json"

	// Логи
	LogInvalidContentType = "получен неверный Content-Type: %s"
	ErrInvalidContentType = "неверный Content-Type, ожидается application/json"

	DefaultProtocol = "http"
)
