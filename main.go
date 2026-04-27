package main

import (
	"archive/zip"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/joho/godotenv"
)

type FtpConfig struct {
	Host         string
	UserName     string
	Password     string
	SaveLocation string
	WorldId      string
}

func main() {
	env := os.Getenv("ENV")
	if env != "prod" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("error loading .env file")
		}
	}

	host := os.Getenv("FTP_HOST")
	if host == "" {
		log.Fatal("failed to load FTP_HOST var")
	}

	user := os.Getenv("FTP_USER")
	if user == "" {
		log.Fatal("failed to load FTP_USER var")
	}

	pass := os.Getenv("FTP_PASS")
	if pass == "" {
		log.Fatal("failed to load FTP_PASS var")
	}

	savePath := os.Getenv("WORLD_SAVE_LOCATION")
	if savePath == "" {
		log.Fatal("failed to load WORLD_SAVE_LOCATION var")
	}

	worldId := os.Getenv("WORLD_ID")
	if worldId == "" {
		log.Fatal("failed to load WORLD_ID var")
	}

	config := &FtpConfig{
		Host:         host,
		UserName:     user,
		Password:     pass,
		SaveLocation: savePath,
		WorldId:      worldId,
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
	worldPath := fmt.Sprintf("%s/%s", config.SaveLocation, config.WorldId)
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
