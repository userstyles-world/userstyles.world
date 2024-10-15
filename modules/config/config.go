// Package config provides configuration options.
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type (
	AppConfig struct {
		Debug      bool
		Profiling  bool
		Production bool

		Addr        string
		BaseURL     string
		MonitorURL  string
		ProxyHeader string

		Name         string
		Description  string
		Codename     string
		Copyright    string
		Started      time.Time `json:"-"`
		EmailRe      string
		PageMaxItems int
		CodeMaxItems int

		// External links.
		Discord        string
		Matrix         string
		OpenCollective string
		GitRepository  string

		// Git info.
		GitCommitURL string `json:"-"`
		GitCommitSHA string `json:"-"`
		GitSignature string `json:"-"`
	}

	DatabaseConfig struct {
		Name             string
		Debug            string
		Colorful         bool
		Migrate          bool
		Drop             bool
		RandomData       bool
		RandomDataAmount int
		MaxOpenConns     int
	}

	OpenAuthConfig struct {
		GitHubID       string
		GitHubSecret   string
		GitLabID       string
		GitLabSecret   string
		CodebergID     string
		CodebergSecret string
	}

	SecretsConfig struct {
		PasswordCost           int
		ScrambleStepSize       int
		ScrambleBytesPerInsert int

		SessionTokenKey  string
		RecoverTokenKey  string
		ProviderTokenKey string

		CryptoKey        string
		StatsKey         string
		OAuthClientKey   string
		OAuthProviderKey string
		EmailAddress     string
		EmailPassword    string
		EmailServer      string
	}

	StorageConfig struct {
		DataDir   string
		CacheDir  string
		ImageDir  string
		StyleDir  string
		ProxyDir  string
		PublicDir string
		LogFile   string
	}

	config struct {
		App      AppConfig
		Database DatabaseConfig
		OpenAuth OpenAuthConfig
		Secrets  SecretsConfig
		Storage  StorageConfig
	}
)

var (
	// GoVersion is the version of Go compiler used to build this program.
	GoVersion string

	// GitCommit is the latest Git commit used to build this program.
	GitCommit string

	// GitSignature is the Git version string used to build this program.
	GitSignature string

	// App stores general configuration.
	App *AppConfig

	// Database stores database configuration.
	Database *DatabaseConfig

	// OpenAuth stores configuration needed for connecting to external services.
	OpenAuth *OpenAuthConfig

	// Secrets stores cryptographic keys and related configuration.
	Secrets *SecretsConfig

	// Storage stores paths to directories and files.
	Storage *StorageConfig
)

// UpdateGitInfo updates copyright year at the start of every year.
func (app *AppConfig) UpdateCopyright() {
	if app.Name == "" {
		log.Fatal("config: App.Name can't be an empty string")
	}

	app.Copyright = fmt.Sprintf("© 2020–%d %s", time.Now().Year(), app.Name)
}

// UpdateGitInfo updates dynamic Git-specific fields.
func (app *AppConfig) UpdateGitInfo() {
	if app.GitRepository == "" {
		log.Fatal("config: App.GitRepository can't be an empty string")
	}

	app.GitCommitURL = fmt.Sprintf("%s/commit/%s", app.GitRepository, GitCommit)
	app.GitCommitSHA = fmt.Sprintf("%.8s", GitCommit)
	app.GitSignature = GitSignature
}

// DefaultConfig is a set of default configurations used for development and
// testing environments.
func DefaultConfig() *config {
	return &config{
		App: AppConfig{
			Addr:         ":3000",
			BaseURL:      "http://localhost:3000",
			Name:         "UserStyles.world",
			Description:  "A free and open-source, community-driven website for browsing and sharing UserCSS userstyles.",
			Codename:     "Fennec Fox",
			Started:      time.Now(),
			EmailRe:      `^[a-zA-Z0-9.!#$%&’*+/=?^_\x60{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)+$`,
			PageMaxItems: 36,
			CodeMaxItems: 250,

			Discord:        "https://discord.gg/P5zra4nFS2",
			Matrix:         "https://matrix.to/#/#userstyles:matrix.org",
			OpenCollective: "https://opencollective.com/userstyles",
			GitRepository:  "https://github.com/userstyles-world/userstyles.world",
		},
		Database: DatabaseConfig{
			Name:         "dev.db",
			Debug:        "info",
			Colorful:     true,
			MaxOpenConns: 10,
		},
		Secrets: SecretsConfig{
			PasswordCost:           10,
			ScrambleStepSize:       2,
			ScrambleBytesPerInsert: 3,
			SessionTokenKey:        "ABigSecretPassword",
			RecoverTokenKey:        "OhNoWeCantUseTheSameAsJWTBeCaUseSeCuRiTy1337",
			ProviderTokenKey:       "ImNotACatButILikeUnicorns",
			CryptoKey:              "ABigSecretPasswordWhichIsExact32",
			StatsKey:               "KeyUsedForHashingStatistics",
			OAuthClientKey:         "AnotherStringLetstrySomethΦΦΦ",
			OAuthProviderKey:       "(✿◠‿◠＾◡＾)っ✂❤",
		},
		Storage: StorageConfig{
			DataDir:   "data",
			CacheDir:  "data/cache",
			ImageDir:  "data/images",
			StyleDir:  "data/styles",
			ProxyDir:  "data/proxy",
			PublicDir: "data/public",
			LogFile:   "data/userstyles.log",
		},
	}
}

// Load tries to load configuration from a given path.
func Load(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	c := DefaultConfig()
	if err = json.Unmarshal(b, &c); err != nil {
		return err
	}

	c.App.UpdateGitInfo()
	c.App.UpdateCopyright()

	if c.App.Debug {
		b, err := json.MarshalIndent(c, "", "\t")
		if err != nil {
			return err
		}

		log.Println("config:", string(b))
	}

	App = &c.App
	Database = &c.Database
	OpenAuth = &c.OpenAuth
	Secrets = &c.Secrets
	Storage = &c.Storage

	return nil
}
