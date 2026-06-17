package main

import (
	"encoding/json"
	"io"
	"os"

	templates "bartoostveen.nl/bulkmailer"
	"bartoostveen.nl/bulkmailer/config"
	"bartoostveen.nl/bulkmailer/job"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, PadLevelText: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	err, cfg := config.Load()
	if err != nil {
		log.Fatal("Failed to load app config", err)
	}

	_ = os.Mkdir(cfg.TargetDir, 0755)
	jobsFile, err := os.Open(cfg.JobsFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(jobsFile *os.File) {
		_ = jobsFile.Close()
	}(jobsFile)

	jobsText, err := io.ReadAll(jobsFile)
	if err != nil {
		log.Fatalln(err)
	}

	var jobs []job.EmailJob
	err = json.Unmarshal(jobsText, &jobs)
	if err != nil {
		log.Fatalln(err)
	}

	job.ProcessAllJobs(cfg, templates.Templates, jobs)
}
