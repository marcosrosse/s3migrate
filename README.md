# About this repository

This was a test developed for a SRE job opportunity.

## Legacy asset move

Let's say we have backend service whose mission is to store users avatars (in PNG files) by uploading them to S3 and then their S3 paths (without the domain part that shows the S3 bucket name) is eventually stored in a SQL database table for future reference. These paths are used to show these PNG files in the frontend.

Let's say that for historical reasons, in the past, all these PNG files were stored in one bucket with one path (directory structure) (from now on `legacy-s3`) and after some time the bucket changed and the path within the bucket changed as well (from now on `production-s3`). Which means, half of the PNGs live now inside one bucket with one path and the other half live in a different bucket with a different path. And, of course, this means that the database table contains a mix of paths.

Example:

- One legacy PNG would have a URL like https://legacy-url/image/avatar-32425.png. In this case the database would contain `image/avatar-32425.png`
- One modern PNG would have a URL like https://modern-url/avatar/avatar-32425.png. In this case the database would contain `avatar/avatar-32425.png`

### Your task

You need to write a program that moves all images from the `legacy-s3` to the `production-s3` and updates their paths in the database. To clarify, the program will make sure that all objects in the legacy bucket/path are correctly moved to the new one. This means, that at the end of the execution, the database will also contain only paths with the modern prefix.

# How to run

### Prepare the environment
Create the Postgres and MinIO containers.

podman-compose up -d

Execute the script seed.py + quantity of files to be created in the bucket and database

### Build the project
``go build -o s3migrate cmd/main.go``

### Export variables
Configure the Postgres and Minio/S3 variables in the file scripts/export-envs.sh and run the command: ``eval "$(scripts/export-envs.sh)"``

### Run the binary
``./s3migrate -h`` for help information.

``./s3migrate -srcBucket=legacy-s3 -dstBucket=production-s3 -srcPath=image/ -dstPath=avatar/ 
``
