package config

import (
	"fmt"
	"path"
	"time"
)

type ScrambleSettings struct {
	StepSize       int
	BytesPerInsert int
}

var (
	GitCommit  string
	GitVersion string

	Port                 = getEnv("PORT", ":3000")
	BaseURL              = getEnv("BASE_URL", "http://localhost"+Port)
	DB                   = getEnv("DB", "dev.db")
	DBDebug              = getEnv("DB_DEBUG", "silent")
	DBColor              = getEnvBool("DB_COLOR", false)
	DBMigrate            = getEnvBool("DB_MIGRATE", false)
	DBDrop               = getEnvBool("DB_DROP", false)
	DBRandomData         = getEnvBool("DB_RANDOM_DATA", false)
	DBRandomDataAmount   = getEnvInt("DB_RANDOM_DATA_AMOUNT", 100)
	Salt                 = getEnvInt("SALT", 10)
	JWTSigningKey        = getEnv("JWT_SIGNING_KEY", "ABigSecretPassword")
	VerifyJWTSigningKey  = getEnv("VERIFY_JWT_SIGNING_KEY", "OhNoWeCantUseTheSameAsJWTBeCaUseSeCuRiTy1337")
	OAuthpJWTSigningKey  = getEnv("OAUTHP_JWT_SIGNING_KEY", "ImNotACatButILikeUnicorns")
	CryptoKey            = getEnv("CRYPTO_KEY", "ABigSecretPasswordWhichIsExact32")
	StatsKey             = getEnv("STATS_KEY", "KeyUsedForHashingStatistics")
	OAuthKey             = getEnv("OAUTH_KEY", "AnotherStringLetstrySomethΦΦΦ")
	OAuthpKey            = getEnv("OAUTHP_KEY", "(✿◠‿◠＾◡＾)っ✂❤")
	EmailAddress         = getEnv("EMAIL_ADDRESS", "test@userstyles.world")
	EmailPassword        = getEnv("EMAIL_PWD", "hahah_not_your_password")
	GitHubClientID       = getEnv("GITHUB_CLIENT_ID", "SOmeOneGiVeMeIdEaSwHaTtOpUtHeRe")
	GitHubClientSecret   = getEnv("GITHUB_CLIENT_SECRET", "OurSecretHere?_www.youtube.com/watch?v=dQw4w9WgXcQ")
	GitlabClientID       = getEnv("GITLAB_CLIENT_ID", "SOmeOneGiVeMeIdEaSwHaTtOpUtHeRe")
	GitlabClientSecret   = getEnv("GITLAB_CLIENT_SECRET", "www.youtube.com/watch?v=dQw4w9WgXcQ")
	CodebergClientID     = getEnv("CODEBERG_CLIENT_ID", "SOmeOneGiVeMeIdEaSwHaTtOpUtHeRe")
	CodebergClientSecret = getEnv("CODEBERG_CLIENT_SECRET", "IMgettinggboredd")
	PerformanceMonitor   = getEnvBool("PERFORMANCE_MONITOR", false)
	IMAPServer           = getEnv("IMAP_SERVER", "mail.userstyles.world:587")
	ProxyMonitor         = getEnv("PROXY_MONITOR", "unset")
	SearchReindex        = getEnvBool("SEARCH_REINDEX", false)

	// Production is used for various "feature flags".
	Production = DB != "dev.db"

	ScrambleConfig = &ScrambleSettings{
		StepSize:       getEnvInt("NONCE_SCRAMBLE_STEP", 2),
		BytesPerInsert: getEnvInt("NONCE_SCRAMBLE_BYTES_PER_INSERT", 3),
	}

	DataDir   = path.Join(getEnv("DATA_DIR", "data"))
	CacheDir  = path.Join(DataDir, "cache")
	ImageDir  = path.Join(DataDir, "images")
	ProxyDir  = path.Join(DataDir, "proxy")
	PublicDir = path.Join(DataDir, "public")

	LogFile = path.Join(DataDir, "userstyles.log")

	AppName         = "UserStyles.world"
	AppCodeName     = "Fennec Fox"
	AppSourceCode   = "https://github.com/userstyles-world/userstyles.world"
	AppLatestCommit = AppSourceCode + "/commit/" + GitCommit
	AppCommitSHA    = fmt.Sprintf("%.7s", GitCommit)
	AppUptime       = time.Now()
	AppPageMaxItems = 36

	AppLinkChatDiscord = "https://discord.gg/P5zra4nFS2"
	AppLinkChatMatrix  = "https://matrix.to/#/#userstyles:matrix.org"

	AllowedEmailsRe = `^[a-zA-Z0-9.!#$%&’*+/=?^_\x60{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)+$`
	AllowedImagesRe = `^https:\/\/(www\.)?(userstyles.world/api|(user-images|raw|gist)\.githubusercontent\.com|github\.com|gitlab\.com|codeberg\.org|cdn\.jsdelivr\.net/gh/33kk)\/.*\.(jpe?g|png|webp)(\?inline=(true|false))?$`

	CachedCodeItems = uint(getEnvInt("CACHED_CODE_ITEMS", 25))
)

// OAuthURL returns the proper callback URL depending on the environment.
func OAuthURL() string {
	return BaseURL + "/api/callback/"
}
