package config

import (
	"flag"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AppConfig struct {
	LogLevel        log.Level  `mapstructure:"log_level"`
	From            string     `mapstructure:"from"`
	ReplyTo         string     `mapstructure:"reply_to"`
	Subject         string     `mapstructure:"subject"`
	TargetDir       string     `mapstructure:"target_dir"`
	SMTP            SMTPConfig `mapstructure:"smtp"`
	Dry             bool       `mapstructure:"dry"`
	JobsFile        string     `mapstructure:"jobs_file"`
	UniqueFileNames bool       `mapstructure:"unique_file_names"`
	Paralellism     int        `mapstructure:"paralellism"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func assertNotBlank(str string, name string) {
	if strings.TrimSpace(str) == "" {
		log.Fatalf("%s must not be blank!", name)
	}
}

func (c *AppConfig) validate() {
	if !c.Dry {
		assertNotBlank(c.SMTP.Host, "SMTP.Host")
	}

	assertNotBlank(c.From, "From")
	assertNotBlank(c.Subject, "Subject")
	assertNotBlank(c.TargetDir, "TargetDir")
	assertNotBlank(c.JobsFile, "JobsFile")
}

func Load() (error, AppConfig) {
	var cfg AppConfig

	configPath := flag.String("c", "", "path to config file")
	flag.Parse()

	v := viper.New()

	v.SetDefault("log_level", log.InfoLevel)
	v.SetDefault("dry", true)
	v.SetDefault("target_dir", "target")
	v.SetDefault("jobs_file", "jobs.json")
	v.SetDefault("unique_file_names", false)
	v.SetDefault("paralellism", -1)

	if *configPath != "" {
		v.SetConfigFile(*configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
	}

	if err := v.ReadInConfig(); err != nil {
		// Only error if explicitly specified
		if *configPath != "" {
			return err, cfg
		}

		log.WithError(err).Warn("Could not load config file, silently failing")
	}

	_ = v.BindEnv("from")
	_ = v.BindEnv("reply_to")
	_ = v.BindEnv("subject")
	_ = v.BindEnv("target_dir")

	_ = v.BindEnv("smtp.host")
	_ = v.BindEnv("smtp.port")
	_ = v.BindEnv("smtp.username")
	_ = v.BindEnv("smtp.password")

	_ = v.BindEnv("jobs_file")
	_ = v.BindEnv("unique_file_names")
	_ = v.BindEnv("dry")

	v.SetEnvPrefix("BULKMAILER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.AutomaticEnv()

	if err := v.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unable to decode config: %w", err), cfg
	}

	log.SetLevel(cfg.LogLevel)

	cfg.validate()

	log.Tracef("Configuration:\n%+v", cfg)

	return nil, cfg
}
