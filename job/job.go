package job

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"text/template"
	"time"

	"bartoostveen.nl/bulkmailer/config"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/go-mail"
)

type EmailJob struct {
	Recipient string      `json:"recipient"`
	Cc        []string    `json:"cc,omitempty"`
	Bcc       []string    `json:"bcc,omitempty"`
	Template  string      `json:"template"`
	Data      interface{} `json:"extraAttrs"`
}

func ProcessJob(cfg config.AppConfig, templates *template.Template, job EmailJob) error {
	var body bytes.Buffer
	templateName := job.Template
	if strings.TrimSpace(templateName) == "" {
		templateName = "template.txt"
	}

	if err := templates.ExecuteTemplate(&body, templateName, job); err != nil {
		log.Tracef("Error rendering template, job was: %+v", job)
		return err
	}

	m := mail.NewMsg()
	if err := m.From(cfg.From); err != nil {
		return err
	}

	if err := m.To(job.Recipient); err != nil {
		return err
	}

	if len(job.Cc) != 0 {
		if err := m.Cc(job.Cc...); err != nil {
			return err
		}
	}
	if len(job.Bcc) != 0 {
		if err := m.Bcc(job.Bcc...); err != nil {
			return err
		}
	}

	if cfg.ReplyTo != "" {
		if err := m.ReplyTo(cfg.ReplyTo); err != nil {
			return err
		}
	}

	m.Subject(cfg.Subject)
	m.SetBodyString(mail.TypeTextPlain, body.String())

	id, err := generateMessageID(cfg.SMTP.Host, 22)
	if err != nil {
		return err
	}
	m.SetMessageIDWithValue(id)

	var fileName string
	escapedRecipient := strings.ReplaceAll(job.Recipient, "/", "_")
	if cfg.UniqueFileNames {
		fileName = fmt.Sprintf("%s/%s-%d.eml", cfg.TargetDir, escapedRecipient, time.Now().UnixMilli())
	} else {
		fileName = fmt.Sprintf("%s/%s.eml", cfg.TargetDir, escapedRecipient)
	}

	if err := m.WriteToFile(fileName); err != nil {
		return err
	}

	if cfg.Dry {
		return nil
	}

	var c *mail.Client
	if cfg.SMTP.Username == "" || cfg.SMTP.Password == "" {
		c, err = mail.NewClient(
			cfg.SMTP.Host,
			mail.WithPort(cfg.SMTP.Port),
		)
	} else {
		c, err = mail.NewClient(
			cfg.SMTP.Host,
			mail.WithPort(cfg.SMTP.Port),
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(cfg.SMTP.Username),
			mail.WithPassword(cfg.SMTP.Password),
		)
	}
	if err != nil {
		return err
	}

	for try := range cfg.Retries + 1 {
		err = c.DialAndSend(m)
		if err == nil || try >= cfg.Retries {
			return err
		}
		delay := 250 * (1 << try)
		log.WithError(err).Warnf("Retrying sending email to %s after %d (attempt %d/%d)", job.Recipient, delay, try, cfg.Retries)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	return nil // unreachable
}

func ProcessAllJobs(cfg config.AppConfig, templates *template.Template, jobs []EmailJob) {
	var wg sync.WaitGroup
	ch := make(chan EmailJob)

	var jobCount int
	if cfg.Paralellism == -1 {
		jobCount = runtime.NumCPU()
	} else {
		jobCount = max(1, min(cfg.Paralellism, runtime.NumCPU()))
	}

	for cpu := range jobCount {
		wg.Add(1)

		go (func() {
			log.Infof("Starting thread t%d...\n", cpu)
			defer wg.Done()

			for j := range ch {
				if cfg.Dry {
					log.Infof("[t%d] saving draft email for %s\n", cpu, j.Recipient)
				} else {
					log.Infof("[t%d] saving and sending email to %s", cpu, j.Recipient)
				}
				err := ProcessJob(cfg, templates, j)
				if err != nil {
					log.Warnf("[t%d] [%s] Error: %v\n", cpu, j.Recipient, err)
				}
			}
		})()
	}
	for i := range jobs {
		ch <- jobs[i]
	}
	close(ch)
	wg.Wait()
}

func generateMessageID(domain string, length int) (string, error) {
	rnd := make([]byte, length)
	n, err := rand.Read(rnd) // CSPRNG
	if err != nil || n < length {
		if err == nil {
			err = fmt.Errorf("failed to generate %d random bytes", n)
		}

		return "", err
	}

	id := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(rnd)
	return fmt.Sprintf("%s@%s", id, domain), nil
}
