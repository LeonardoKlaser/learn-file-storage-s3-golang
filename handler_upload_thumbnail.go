package main

import (
	"fmt"
	"net/http"
	"io"
	"log"
	"encoding/base64"
	"crypto/rand"
	"strings"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
	"os"
	"mime"
	"errors"
	"path/filepath"
	
)

func getThumbnailURL(root,extension, videoID string) string {
	path := filepath.Join(root, fmt.Sprintf("%s.%s", videoID, extension))
	return path
}

func getFileExtension(contentType string) (string, error) {
	splited := strings.Split(contentType, "/")
	if len(splited) < 2 {
		return "", errors.New("Error to get file extension") 
	}
	return splited[1], nil
}

func ValidateMediaType(media string) error {
	mediaType, _ , err := mime.ParseMediaType(media)
	if err != nil{
		return errors.New("Invalid Media Type")
	}

	if mediaType != "image/jpeg" && mediaType != "image/png" {
		 return errors.New("Invalid Media Type")
	}

	return nil
}

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}


	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	const maxMemory = 10 << 20

	r.ParseMultipartForm(maxMemory)
	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form File", err)
		return
	}

	contentType := header.Header.Get("Content-Type")

	err = ValidateMediaType(contentType)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error to read media Type: ", err)
                return
	}
	
	video, err := cfg.db.GetVideo(videoID)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Error to get video from database" ,err)
		return
	}
	log.Println(file)
	if video.UserID.String() != userID.String(){
		respondWithError(w, http.StatusUnauthorized, "The authenticated user is not the owner of the video", err)
		return
	}

	extension, err := getFileExtension(contentType)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Error to get file extension", err)
                return
	}

	key := make([]byte, 32)
	rand.Read(key)
	encodedString := base64.RawURLEncoding.EncodeToString(key)

	path := getThumbnailURL(cfg.assetsRoot, extension, encodedString)
	URL := fmt.Sprintf("http://localhost:%s/assets/%s.%s", cfg.port, encodedString, extension)
	video.ThumbnailURL = &URL
	
	imageFile, err := os.Create(path)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Error to craete new file with image", err)
                return
        }
	
	_, err = io.Copy(imageFile, file)

	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Error to insert image metadata to assets file", err)
                return
        }

	err = cfg.db.UpdateVideo(video)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error to update Video" ,err)
		return
	}


	respondWithJSON(w, http.StatusOK, video)
}
