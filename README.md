# About this repository

This was a test developed for a SRE job opportunity.

# Setup used
Components and versions:

- SO: Fedora Linux 38 Workstation
- Go version 1.21.0
- Podman version 4.7.0
- MinIO container latest version
- PostgreSQL container latest version

I created the Postgres and MinIO containers to emulate the AWS S3 Bucket and Production database.

# How to run

### Prepare the environment
To create the Postgres and MinIO containers, run access the repository root path and run:

``compose up -d``

Execute the script seed.py + quantity of files to be created in the bucket and database

### Build the project
Source code:
``go build -o s3migrate cmd/main.go``

Container image:
``podman build -t <REGISTRY URL>/<IMAGE NAME>:<VERSION> -f Dockerfile``

### Export variables
Configure the Postgres and Minio/S3 variables with the user, password, host, etc, in the file scripts/export-envs.sh and run the command: ``eval "$(scripts/export-envs.sh)"``

### Run CLI App
#### Binary file:
After export all variables with the script export-envs.sh running the binary:
``./s3migrate``
#### Container image:
``podman run -d --network=s3migrate_cli <REGISTRY URL>/<IMAGE NAME>:<VERSION> /s3migrate``
