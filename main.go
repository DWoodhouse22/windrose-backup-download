package main

import (
	"archive/zip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

type FtpConfig struct {
	Host     string `json:"host"`
	UserName string `json:"username"`
	Password string `json:"password"`
	WorldId  string `json:"world_id"`
}

const saveDir = "windrose/R5/Saved/SaveProfiles/Default/RocksDB/0.10.0/Worlds"

func main() {
	configPath := "config.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	c, err := ftpConnect(config)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	files, err := listSaveFiles(c, config)
	if err != nil {
		log.Fatal(err)
	}

	zipPath, err := saveFiles(c, config, files)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("created zip:", zipPath)
}

func loadConfig(path string) (*FtpConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg FtpConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if cfg.Host == "" || cfg.Password == "" || cfg.UserName == "" || cfg.WorldId == "" {
		return nil, fmt.Errorf("invalid config file")
	}
	return &cfg, nil
}

func ftpConnect(config *FtpConfig) (*ftp.ServerConn, error) {
	c, err := ftp.Dial(
		config.Host,
		ftp.DialWithTimeout(5*time.Second),
		ftp.DialWithExplicitTLS(&tls.Config{
			InsecureSkipVerify: true,
		}),
	)
	if err != nil {
		return nil, err
	}

	if err := c.Login(config.UserName, config.Password); err != nil {
		return nil, err
	}

	return c, nil
}

func listSaveFiles(c *ftp.ServerConn, config *FtpConfig) ([]string, error) {
	worldPath := fmt.Sprintf("%s/%s", saveDir, config.WorldId)
	files, err := c.NameList(worldPath)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func saveFiles(c *ftp.ServerConn, config *FtpConfig, paths []string) (string, error) {
	zipFileName := fmt.Sprintf("%s_%s.zip", config.WorldId, timestamp())
	outFile, err := os.Create(zipFileName)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	for _, p := range paths {
		r, err := c.Retr(p)
		if err != nil {
			return "", err
		}

		fileName := strings.TrimPrefix(p, "/")
		fmt.Printf("downloading %s...\n", fileName)
		w, err := zipWriter.Create(fileName)
		if err != nil {
			r.Close()
			return "", err
		}

		if _, err := io.Copy(w, r); err != nil {
			r.Close()
			return "", err
		}
		r.Close()
	}

	return zipFileName, nil
}

func timestamp() string {
	return time.Now().Format("2006-01-02_15-04-05")
}
