package main

import (
	"fmt"
	"net/http"
	"io"
	"log"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func getThumbnailURL(port, videoID, fileExtension string) *string {
	thumbnailURL := fmt.Sprintf("http://localhost:%s/api/thumbnails/%s.%s", port, videoID, fileExtension)
	return &thumbnailURL
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

	b, err := io.ReadAll(file)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error to read the file", err)
		return
	}
	
	video, err := cfg.db.GetVideo(videoID)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Error to get video from database" ,err)
		return
	}
	log.Println(video)
	if video.UserID.String() != userID.String(){
		respondWithError(w, http.StatusUnauthorized, "The authenticated user is not the owner of the video", err)
		return
	}

	videoThumb := thumbnail{
		data : b,
		mediaType : contentType,
	}

	videoThumbnails[videoID] =  videoThumb
	video.ThumbnailURL = getThumbnailURL(cfg.port, video.ID.String(), "png")
	log.Println(video)
	err = cfg.db.UpdateVideo(video)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error to update Video" ,err)
		return
	}


	respondWithJSON(w, http.StatusOK, video)
}
