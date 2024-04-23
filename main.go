package main

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

func main() {
    router := SetupRouter()
    router.Run(":8080")
}

func SetupRouter() *gin.Engine {
    router := gin.Default()

    router.POST("/upload", UploadImage)
    router.GET("/image/:imageID/:filetype", GetImage)
    router.GET("/transform/rotate/:imageID/:angle", RotateImage)
    router.GET("/transform/resize/:imageID/:width/:height", ResizeImage)

    return router
}

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a unique filename without the extension
	filename := uuid.New().String()
	if err := c.SaveUploadedFile(file, "./uploads/"+filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": filename})
}

// Define wrapper functions for encoders
func encodeJPEG(w io.Writer, m image.Image) error {
	return jpeg.Encode(w, m, nil)
}

func encodeGIF(w io.Writer, m image.Image) error {
	return gif.Encode(w, m, nil)
}

func encodePNG(w io.Writer, m image.Image) error {
	return png.Encode(w, m)
}

// Define map for encoders
var encoders = map[string]func(io.Writer, image.Image) error{
	"jpeg": encodeJPEG,
	"jpg": encodeJPEG,
	"gif":  encodeGIF,
	"png":  encodePNG,
}

func GetImage(c *gin.Context) {
	imageID := c.Param("imageID")
	imageType := c.Param("filetype") // File type provided as parameter

	imagePath := "./uploads/" + imageID

	// Open the image file
	imageFile, err := os.Open(imagePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}
	defer imageFile.Close()

	// Decode the image to image.Image
	img, _, err := image.Decode(imageFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode image"})
		return
	}

	// Create a buffer to write the image
	var imageBuffer bytes.Buffer


	// Check if the encoding function exists for the requested image type
	encodeFunc, ok := encoders[imageType]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type"})
		return
	}

	// Encode the image to the desired format
	if err := encodeFunc(&imageBuffer, img); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode image"})
		return
	}

	// Set the Content-Type header based on the requested file type
	contentType := "image/" + imageType

	// Set the Content-Type header in the response
	c.Header("Content-Type", contentType)

	// Write the image buffer to the response
	c.Data(http.StatusOK, contentType, imageBuffer.Bytes())
}

func RotateImage(c *gin.Context) {
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

	img, format, err := image.Decode(imageFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode image"})
		return
	}

	rotatedImage := rotate(img, angle)

	rotatedFilename := fmt.Sprintf("%s_rotated_%s.%s", strings.TrimSuffix(imageID, filepath.Ext(imageID)), angleStr, format)
	rotatedPath := "./uploads/" + rotatedFilename
	out, err := os.Create(rotatedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save rotated image"})
		return
	}
	defer out.Close()

	switch format {
	case "jpeg":
		if err := jpeg.Encode(out, rotatedImage, nil); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode rotated image"})
			return
		}
	case "png":
		if err := png.Encode(out, rotatedImage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode rotated image"})
			return
		}
	// Add cases for other image formats as needed
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unsupported image format"})
		return
	}

	c.File(rotatedPath)
}


func rotate(img image.Image, angle float64) image.Image {
    bounds := img.Bounds()
    rotated := image.NewRGBA(bounds)

    // Calculate the midpoints
    midX := float64(bounds.Max.X) / 2
    midY := float64(bounds.Max.Y) / 2

    // Convert degrees to radians
    angleRad := angle * math.Pi / 180

    for x := bounds.Min.X; x < bounds.Max.X; x++ {
        for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
            // Subtract midpoints to center the rotation around the origin
            px := float64(x) - midX
            py := float64(y) - midY

            // Perform rotation
            newX := px*math.Cos(angleRad) - py*math.Sin(angleRad)
            newY := px*math.Sin(angleRad) + py*math.Cos(angleRad)

            // Translate the image back to its original position
            newX += midX
            newY += midY

            // Convert float coordinates to integer
            intX := int(newX + 0.5)
            intY := int(newY + 0.5)

            // Check if the new coordinates are within bounds
            if intX >= bounds.Min.X && intX < bounds.Max.X && intY >= bounds.Min.Y && intY < bounds.Max.Y {
                // Set the pixel color at the new coordinates
                rotated.Set(x, y, img.At(intX, intY))
            }
        }
    }

    return rotated
}

func ResizeImage(c *gin.Context) {
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

	img, format, err := image.Decode(imageFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode image"})
		return
	}

	resizedImage := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	resizedFilename := fmt.Sprintf("%s_resized_%dx%d.%s", strings.TrimSuffix(imageID, filepath.Ext(imageID)), width, height, format)
	resizedPath := "./uploads/" + resizedFilename
	out, err := os.Create(resizedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resized image"})
		return
	}
	defer out.Close()

	switch format {
	case "jpeg":
		if err := jpeg.Encode(out, resizedImage, nil); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode resized image"})
			return
		}
	case "png":
		if err := png.Encode(out, resizedImage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode resized image"})
			return
		}
	// Add cases for other image formats as needed
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unsupported image format"})
		return
	}

	c.File(resizedPath)
}

