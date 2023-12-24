package normal

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func DownloadToLocal(url, localPath string, wg *sync.WaitGroup) {
	defer wg.Done()

	response, err := http.Get(url)
	if err != nil {
		log.Printf("Error downloading image from %s: %v", url, err)
		return
	}
	defer response.Body.Close()

	file, err := os.Create(localPath)
	if err != nil {
		log.Printf("Error creating file %s: %v", localPath, err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Printf("Error writing to file %s: %v", localPath, err)
		return
	}

	log.Printf("Downloaded %s to %s", url, localPath)
}
