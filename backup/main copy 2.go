package main

import (
	"fmt"
	"math"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	router := gin.Default()

	router.POST("/upload", uploadImage)
	router.GET("/image/:imageID", getImage)
	router.GET("/transform/rotate/:imageID/:angle", rotateImage)
	router.GET("/transform/resize/:imageID/:width/:height", resizeImage)

	router.Run(":8080")
}

func uploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(file.Filename))
	if err := c.SaveUploadedFile(file, "./uploads/"+filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": filename})
}

func getImage(c *gin.Context) {
	imageID := c.Param("imageID")
	c.File("./uploads/" + imageID)
}

func rotateImage(c *gin.Context) {
	imageID := c.Param("imageID")
	angleStr := c.Param("angle")

	angle, err := strconv.ParseFloat(angleStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rotation angle"})
		return
	}

	imagePath := "./uploads/" + imageID
	imageFile, err := os.Open(imagePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}
	defer imageFile.Close()

	img, _, err := image.Decode(imageFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode image"})
		return
	}

	rotatedImage := rotate(img, angle)
	rotatedFilename := fmt.Sprintf("%s_rotated_%s.jpg", strings.TrimSuffix(imageID, filepath.Ext(imageID)), angleStr)
	rotatedPath := "./uploads/" + rotatedFilename
	out, err := os.Create(rotatedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save rotated image"})
		return
	}
	defer out.Close()

	if err := jpeg.Encode(out, rotatedImage, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode rotated image"})
		return
	}

	c.File(rotatedPath)
}

func rotate(img image.Image, angle float64) image.Image {
	// Create a new image with the same bounds as the original image
	b := img.Bounds()
	rotated := image.NewRGBA(b)

	// Define the rotation matrix
	theta := angle * math.Pi / 180
	cos := math.Cos(theta)
	sin := math.Sin(theta)

	// Perform the rotation
	for x := 0; x < b.Max.X; x++ {
		for y := 0; y < b.Max.Y; y++ {
			// Translate pixel so center is at origin
			xp := float64(x) - float64(b.Max.X)/2
			yp := float64(y) - float64(b.Max.Y)/2

			// Rotate pixel
			xr := xp*cos - yp*sin
			yr := xp*sin + yp*cos

			// Translate pixel back to original position
			xp = xr + float64(b.Max.X)/2
			yp = yr + float64(b.Max.Y)/2

			// Set pixel in rotated image
			if xp >= 0 && xp < float64(b.Max.X) && yp >= 0 && yp < float64(b.Max.Y) {
				rotated.Set(int(x), int(y), img.At(int(xp), int(yp)))
			}
		}
	}

	return rotated
}


func resizeImage(c *gin.Context) {
	imageID := c.Param("imageID")
	widthStr := c.Param("width")
	heightStr := c.Param("height")

	width, err := strconv.Atoi(widthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid width"})
		return
	}

	height, err := strconv.Atoi(heightStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid height"})
		return
	}

	imagePath := "./uploads/" + imageID
	imageFile, err := os.Open(imagePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}
	defer imageFile.Close()

	img, _, err := image.Decode(imageFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode image"})
		return
	}

	resizedImage := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	resizedFilename := fmt.Sprintf("%s_resized_%dx%d.jpg", strings.TrimSuffix(imageID, filepath.Ext(imageID)), width, height)
	resizedPath := "./uploads/" + resizedFilename
	out, err := os.Create(resizedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resized image"})
		return
	}
	defer out.Close()

	if err := jpeg.Encode(out, resizedImage, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode resized image"})
		return
	}

	c.File(resizedPath)
}
