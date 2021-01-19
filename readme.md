# Memrizr

This is the repository for the Memrizr memorization application tutorial which is a work in progress. If you're interested in the tutorial, check out the videos on [YouTube](https://www.youtube.com/playlist?list=PLnrGn4P6C4P6yasdEJnEUhueTjCGXGuFe)!

You can find my written tutorials on my [dev.to feed](https://dev.to/jacobsngoodwin).

For the complete "prototype" of this application (I'm making several changes in this tutorial), check out the [wordmem repo](https://github.com/JacobSNGoodwin/wordmem).

## Application Overview

A chart of the tools and applications used in this tutorial is given below.

![App Overview](./application_overview.png)

## Runing the Application

`.env.dev` files are provided in each application directory for providing development environment variables and quickly running this application, hopefully without any major headaches. 

We won't be adding any critical keys directly to the `.env.dev` file. However, we will eventually refer to access key files in `.env.dev`. Make sure to add these key files to your .gitignore.

### Add Test Domain to Host File

This application uses Traefik as a reverse proxy. In `docker-compose.yml`, our `reverse-proxy` Traefik service configures HTTP routes with HOST `malcorp.test`. Therefore, you will need to add the following to your hosts file to map this domain name to localhost.

`127.0.0.1       malcorp.test`

The hosts file can be found (to the best of my knowledge) at `~/etc/hosts` on Mac/Linux or `C:\Windows\System32\drivers\etc` on Windows. 

### Run make init

I have created a script in `Makefile` which creates necessary keypairs (for creating dev and test JSON Web Tokens), runs any docker-compose service for Postgres, migrates our database tables, then shuts down docker-compose. I may add more initializatin commands later on.

As an example, we'll migrate database changes found in `~/account/migrations` to our `postgres-account` service found in `~/docker-compose.yml`.

From the project root director, run:

```bash
make init
```

*Note that in `make init` I do not check to make sure Postgres is ready for connections, even though the docker-compose command will have completed. I want to avoid writting a [complex script](https://stackoverflow.com/questions/57514720/bash-script-command-to-wait-until-docker-compose-process-has-finished-before-mov). Thereofore, there is a chance that the database will not be ready in time for the migrate commands. In msot cases, you just need to try running `make init` again.* 

*If any of you has scripting skills, please submit a PR and I'll update this! It make require pinging, or adding some script into, the Postgres container(s)*

### Add Cloud Key

In order to access Google Cloud for storing profile images, you will need to download a service account JSON file to your `account` application folder and call it `serviceAccount.json`. This file will be references in .env.dev.

Instructions for installing the Google Cloud Storage Client and getting this key are found at:

[https://cloud.google.com/storage/docs/reference/libraries](https://cloud.google.com/storage/docs/reference/libraries)

### Run

To run this code, you will need docker and docker-compose installed on your machine. In the project root, run `docker-compose up`.


Cheers, eh!
