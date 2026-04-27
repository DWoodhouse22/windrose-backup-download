package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("ENV")
	if env != "prod" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("error loading .env file")
		}
	}

	c, err := ftpConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	files, err := listSaveFiles(c)
	if err != nil {
		log.Fatal(err)
	}

	zipPath, err := saveFiles(c, files)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("created zip:", zipPath)
}

func ftpConnect() (*ftp.ServerConn, error) {
	host := os.Getenv("FTP_HOST")
	if host == "" {
		return nil, fmt.Errorf("failed to load FTP_HOST var")
	}

	user := os.Getenv("FTP_USER")
	if user == "" {
		return nil, fmt.Errorf("failed to load FTP_USER var")
	}

	pass := os.Getenv("FTP_PASS")
	if pass == "" {
		return nil, fmt.Errorf("failed to load FTP_PASS var")
	}

	c, err := ftp.Dial("ukln082.gamedata.io:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	if err := c.Login(user, pass); err != nil {
		return nil, err
	}

	return c, nil
}

func listSaveFiles(c *ftp.ServerConn) ([]string, error) {
	savePath := os.Getenv("WORLD_SAVE_LOCATION")
	if savePath == "" {
		return nil, fmt.Errorf("failed to load WORLD_SAVE_LOCATION var")
	}
	worldId := os.Getenv("WORLD_ID")
	if worldId == "" {
		return nil, fmt.Errorf("failed to load WORLD_ID var")
	}
	worldPath := fmt.Sprintf("%s/%s", savePath, worldId)

	files, err := c.NameList(worldPath)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func saveFiles(c *ftp.ServerConn, paths []string) (string, error) {
	worldId := os.Getenv("WORLD_ID")
	if worldId == "" {
		return "", fmt.Errorf("failed to load WORLD_ID var")
	}

	zipFileName := fmt.Sprintf("%s_%s.zip", worldId, timestamp())
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

		fileName := filepath.Base(p)
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
