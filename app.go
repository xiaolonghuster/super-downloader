package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"super-downloader/downloader/normal"
	"super-downloader/repository/mongos"
	"super-downloader/repository/mongos/entity"
	"sync"
	"time"
)

const concurrency = 50

func Init() {

}

func main() {
	// Set MongoDB connection information
	mongoURI := "mongodb://localhost:27017"
	dbName := "mongodb://localhost:27017"

	mongoHandler := mongos.Init(mongoURI, dbName)
	defer mongoHandler.Disconnect()

	var err error
	for {
		// Query image information to download
		var images []entity.ImageInfo

		if images, err = mongoHandler.QueryAllImage(dbName, ""); err != nil {
			logrus.Errorf("query images exception:%v", err)
		}

		// If there is not enough data, break out of the loop
		if len(images) == 0 {
			logrus.Infof("No more data in the database. Exiting.")
			break
		}

		// Control concurrency
		var wg sync.WaitGroup
		semaphore := make(chan struct{}, concurrency)

		for _, img := range images {
			wg.Add(1)
			semaphore <- struct{}{} // Occupy semaphore to limit concurrency
			go func(img entity.ImageInfo) {
				defer func() {
					<-semaphore // Release semaphore
				}()

				// Extract file extension from URL
				fileExt := filepath.Ext(img.Url)

				// Generate local file name based on ID and extension
				localFileName := fmt.Sprintf("%s%s", img.ID, fileExt)
				localPath := filepath.Join("download_folder", localFileName)

				normal.DownloadToLocal(img.Url, localPath, &wg)
			}(img)
		}

		// Wait for all downloads to complete
		wg.Wait()

		// Sleep for a short duration before the next iteration
		time.Sleep(time.Second)
	}

	fmt.Println("All downloads completed.")
}
