# Go Image Upload/Download REST API 

This is a simple Go web service for uploading, transforming, and retrieving images.

## Features

- Upload images via HTTP POST request.
- Retrieve uploaded images by their IDs and file types.
- Rotate uploaded images by a specified angle.
- Resize uploaded images to specific dimensions.

## Installation

1. Clone this repository.
2. Navigate to the project directory.
3. Build and run the service.
4. The service will be available at `http://localhost:8080`.

## Usage

### Uploading Images

To upload an image, send a POST request to `/upload` endpoint with the image file in the form data:


### Retrieving Images

To retrieve an uploaded image, send a GET request to `/image/:imageID/:filetype` endpoint, where `:imageID` is the ID of the uploaded image and `:filetype` is the desired file type (jpeg, png, gif):


### Rotating Images

To rotate an uploaded image, send a GET request to `/transform/rotate/:imageID/:angle` endpoint, where `:imageID` is the ID of the uploaded image and `:angle` is the rotation angle in degrees:


### Resizing Images

To resize an uploaded image, send a GET request to `/transform/resize/:imageID/:width/:height` endpoint, where `:imageID` is the ID of the uploaded image, `:width` is the desired width, and `:height` is the desired height:

## Commands

1. Clone
```git clone https://github.com/Varun58/prog-image-service.git ```

2. Build Application
```go build prog-image-service```

4. Run Test Cases
```go test -v```

5. Run Application
```go run main.go```

## Dependencies

- gin-gonic/gin: HTTP web framework for Go
- nfnt/resize: Image resizing library for Go
- google/uuid: UUID generation library for Go

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
