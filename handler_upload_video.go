package main

import (
	"net/http"
	"context"
    	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"encoding/base64"
	"crypto/rand"
	"fmt"
	"log"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/utils"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
	"os"
	"mime"
	"errors"
)

func validateVideoExtension(media string) error{
	mediaType, _, err := mime.ParseMediaType(media)
	if err != nil {
		return errors.New("Couldn´t get media type")
	}
	if mediaType != "video/mp4" {
		return errors.New("Invalid media Type")
	}
	return nil
}


func (cfg *apiConfig) handlerUploadVideo(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn´t findo jwt", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn´t validate jwt", err)
		return
	}


	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}


	const maxMemory = 10 << 30
	r.ParseMultipartForm(maxMemory)

	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error to retrieve the video from database", err)
		return
	}

	if video.UserID.String() != userID.String(){
		respondWithError(w, http.StatusUnauthorized, "The authenticated user isn´t the owner of the video", err)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn´t get the video uploaded", err)
		return 
	}

	defer file.Close()
	
	contentType := header.Header.Get("Content-Type")
	err = validateVideoExtension(contentType)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Media Type", err)
		return
	}

	f, err := os.CreateTemp("", "tempVideoUploaded-*.mp4")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error create new temporary file to storage de video",err)
		return
	}

	defer os.Remove(f.Name()) // clean up

	_, err = io.Copy(f, file)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error to upload video to temporary file", err)
		return
	}

	f.Seek(0, 0)
	

	fileWithFastStart, err := utils.ProcessVideoForFastStart(f.Name())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error to generate a fast start video file", err)
                return
	}
	fileFromTemporary, err := os.Open(fileWithFastStart)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error to get video from temporary file", err)
		return
	}

	defer fileFromTemporary.Close()

	key := make([]byte, 32)
	rand.Read(key)
	log.Println(f.Name())
	ratio, err := utils.GetVideoAspectRatio(f.Name())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error to get the ratio of archive", err)
		return
	}
	encodedString := "/" + ratio + "/" + base64.RawURLEncoding.EncodeToString(key)

	_, err = cfg.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &cfg.s3Bucket,
		Key: &encodedString,
		Body: fileFromTemporary,
		ContentType: aws.String("video/mp4"),})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error to upload video for aws s3", err)
		return
	}

	log.Println(encodedString)
	newUrl := fmt.Sprintf("%s%s.mp4", cfg.s3CfDistribution, encodedString)
	video.VideoURL = &newUrl
	
	err = cfg.db.UpdateVideo(video)
        if err != nil {
                respondWithError(w, http.StatusBadRequest, "Error to update Video" ,err)
                return
        }


	respondWithJSON(w, http.StatusOK, video)



}
