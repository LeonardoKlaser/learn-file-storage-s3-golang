# learn-file-storage-s3-golang-starter (Tubely)

This repo contains the starter code for the Tubely application - the #1 tool for engagement bait - for the "Learn File Servers and CDNs with S3 and CloudFront" [course](https://www.boot.dev/courses/learn-file-servers-s3-cloudfront-golang) on [boot.dev](https://www.boot.dev)

## Quickstart

*This is to be used as a *reference\* in case you need it, you should follow the instructions in the course rather than trying to do everything here.

## 1. Install dependencies

- [Go](https://golang.org/doc/install)
- `go mod download` to download all dependencies
- [FFMPEG](https://ffmpeg.org/download.html) - both `ffmpeg` and `ffprobe` are required to be in your `PATH`.

```bash
# linux
sudo apt update
sudo apt install ffmpeg

# mac
brew update
brew install ffmpeg
```

- [SQLite 3](https://www.sqlite.org/download.html) only required for you to manually inspect the database.

```bash
# linux
sudo apt update
sudo apt install sqlite3

# mac
brew update
brew install sqlite3
```

- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)

## 2. Download sample images and videos

```bash
./samplesdownload.sh
# samples/ dir will be created
# with sample images and videos
```

## 3. Configure environment variables

Copy the `.env.example` file to `.env` and fill in the values.

```bash
cp .env.example .env
```

You'll need to update values in the `.env` file to match your configuration, but _you won't need to do anything here until the course tells you to_.

## 3. Run the server

```bash
go run .
```

- You should see a new database file `tubely.db` created in the root directory.
- You should see a new `assets` directory created in the root directory, this is where the images will be stored.
- You should see a link in your console to open the local web page.



# Tubely - Video Asset Management SaaS

## 🚀 About This Project

**Tubely** is a SaaS (Software as a Service) product designed to help YouTubers efficiently manage their video assets.

The platform will enable users to:
*   Upload and store video files.
*   Serve video content.
*   Add and manage metadata (titles, descriptions, etc.).
*   Version control for video files.
*   Manage thumbnails and other related video metadata.

This project focuses on building a scalable application capable of handling large media files using Go and cloud services.

## 🎯 Project Goals & Key Learnings

The development of Tubely aims to achieve and understand:

*   Effective strategies for handling "large" files (like videos) versus "small" structured data.
*   Building an application using **Go** and **AWS S3** for robust storage and serving of assets.
*   Techniques for managing files both on traditional filesystems and at scale using serverless solutions like AWS S3.
*   Implementation of video streaming, focusing on data usage efficiency and performance.

## ✨ Core Features (Planned)

*   User accounts and authentication.
*   Video file uploads to AWS S3.
*   Secure video file storage and versioning.
*   Streaming of video content.
*   Management of video metadata (titles, descriptions, custom tags).
*   Thumbnail generation and management.

## 🛠️ Technologies Used

*   **Primary Language:** Go
*   **Cloud Storage & Serving:** AWS S3

---

*Project developed by Leonardo Klaser.*
