# About this repository

This was a test developed for a SRE job opportunity.

# Setup used
Components and versions:

- SO: Fedora Linux 38 Workstation
- Go version 1.21.0
- Podman version 4.7.0
- MinIO container latest version
- MinIO library v7
- PostgreSQL container latest version

Postgres and MinIO containers orchestrated to emulate the AWS S3 Bucket and Production database.

The code were written using a method from MinIO to copy object files from a source to a destination, and in the final it update each column in the database.

# How to run

### Prepare the environment
To create the Postgres and MinIO containers, run access the repository root path and run:

``compose up -d``

Execute the script seed.py + quantity of files to be created in the bucket and database.

### Build the project
From source code:
``go build -o s3migrate cmd/main.go``

From container image:
``podman build -t <REGISTRY URL>/<IMAGE NAME>:<VERSION> -f Dockerfile``

### Export variables
Configure the Postgres and Minio/S3 variables with the user, password, host, etc, in the file scripts/export-envs.sh and run the command: ``eval "$(scripts/export-envs.sh)"``

### Run CLI App
#### Binary file:
After export all variables with the script export-envs.sh running the binary:
``./s3migrate``

#### Container image:
``podman run -d --network=s3migrate_cli <REGISTRY URL>/<IMAGE NAME>:<VERSION> /s3migrate``

### Considerations:
This app cli would be deployed in a Kubernetes Job because of it's simplicity and a "one time use" purpose.

Thinking about performance, I would implement a bulk update to update the objects file path in Postgres, because usually databases work better when processing bulk queries.


