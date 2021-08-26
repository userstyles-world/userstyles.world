package config

import (
	"time"
)

type ScrambleSettings struct {
	StepSize       int
	BytesPerInsert int
}

var (
	Port                 = getEnv("PORT", ":3000")
	DB                   = getEnv("DB", "dev.db")
	DBDebug              = getEnv("DB_DEBUG", "silent")
	DBColor              = getEnvBool("DB_COLOR", false)
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
	GitCommit            = getEnv("GIT_COMMIT", "unset")
	GitVersion           = getEnv("GIT_VERSION", "unset")
	GitHubClientID       = getEnv("GITHUB_CLIENT_ID", "SOmeOneGiVeMeIdEaSwHaTtOpUtHeRe")
	GitHubClientSecret   = getEnv("GITHUB_CLIENT_SECRET", "OurSecretHere?_www.youtube.com/watch?v=dQw4w9WgXcQ")
	GitlabClientID       = getEnv("GITLAB_CLIENT_ID", "SOmeOneGiVeMeIdEaSwHaTtOpUtHeRe")
	GitlabClientSecret   = getEnv("GITLAB_CLIENT_SECRET", "www.youtube.com/watch?v=dQw4w9WgXcQ")
	CodebergClientID     = getEnv("CODEBERG_CLIENT_ID", "SOmeOneGiVeMeIdEaSwHaTtOpUtHeRe")
	CodebergClientSecret = getEnv("CODEBERG_CLIENT_SECRET", "IMgettinggboredd")
	PerformanceMonitor   = getEnvBool("PERFORMANCE_MONITOR", false) == true
	IMAPServer           = getEnv("IMAP_SERVER", "mail.userstyles.world:587")

	// Production is used for various "feature flags".
	Production = DB != "dev.db"

	ScrambleConfig = &ScrambleSettings{
		StepSize:       getEnvInt("NONCE_SCRAMBLE_STEP", 2),
		BytesPerInsert: getEnvInt("NONCE_SCRAMBLE_BYTES_PER_INSERT", 3),
	}

	AppName         = "UserStyles.world"
	AppCodeName     = "Silver Fox"
	AppSourceCode   = "https://github.com/userstyles-world/userstyles.world"
	AppLatestCommit = AppSourceCode + "/commit/" + GitCommit
	AppUptime       = time.Now()
	AppPageMaxItems = 40

	AppLinkChatDiscord = "https://discord.gg/P5zra4nFS2"
	AppLinkChatMatrix  = "https://matrix.to/#/#userstyles:matrix.org"
	AppLinkSource      = "https://github.com/userstyles-world/userstyles.world"

	AllowedEmailsRe = `^[a-zA-Z0-9.!#$%&’*+/=?^_\x60{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)+$`
	AllowedImagesRe = `^https:\/\/(www\.)?(userstyles.world/api|(user-images|raw)\.githubusercontent\.com|github\.com|gitlab\.com|codeberg\.org|cdn\.jsdelivr\.net/gh/33kk)\/.*\.(jpe?g|png|webp)(\?inline=(true|false))?$`
)

// OAuthURL returns the proper callback URL depending on the environment.
func OAuthURL() string {
	return BaseURL() + "/api/callback/"
}

// BaseURL returns the proper BaseURL.
func BaseURL() string {
	if Production {
		return "https://userstyles.world"
	}
	return "http://localhost" + Port
}
