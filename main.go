package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	formFieldPatch = "file"
	uploadDirName  = "upload"
)

func createDir(dir string) error {
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
		return nil
	}

	return err
}

func saveFile(src io.ReadSeeker, dst string) error {
	src.Seek(0, 0)

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return err
	}

	return nil
}

func postUpload(c *gin.Context) {
	f, h, err := c.Request.FormFile(formFieldPatch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	pwd, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dir := pwd + "/" + uploadDirName
	if err := createDir(dir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	path := dir + "/" + h.Filename
	if err := saveFile(f, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(path)

	c.AbortWithStatus(http.StatusOK)
}

func main() {
	r := gin.Default()

	r.POST("/upload", postUpload)
	r.Run(":8080")
}
