# GiveGetGo Backend

GiveGetGo backend repo.

## File Structure

```bash
.
├── .dockerignore
├── .env
├── .github
│  ├── CODEOWNERS
│  └── workflows
│     ├── cd.yml
│     ├── ci.yml
│     └── dev_cd.yml
├── argocd
│  ├── development
│  └── production
├── charts
│  └── application template
│     ├── .helmignore
│     ├── Chart.yaml
│     ├── templates
│     └── values.yaml
├── docker-compose.yml
├── encode-envs.sh
├── entrypoint.sh
├── go.work
├── go.work.sum
├── nginx
│  ├── Dockerfile
│  └── nginx.conf
├── README.md
├── redis
│  ├── .env.redis
│  ├── Dockerfile
│  └── entrypoint.sh
└── servers
   └── service template
      ├── .env.service
      ├── controller
      ├── db
      ├── Dockerfile
      ├── main.go
      ├── middleware
      ├── schema
      ├── server
      └── utils
```
Root Directory: Contains configuration files like .dockerignore, .env, and docker-compose.yml, among others. It also includes the main directories for workflows, configurations, and service templates.

argocd Directory: Contains folders for different deployment environments (development and production).

charts Directory: Includes an application template for Helm charts.

nginx Directory: Houses the Dockerfile and configuration file nginx.conf for setting up an Nginx server.

redis Directory: Contains Redis-specific files such as .env.redis, Dockerfile, and entrypoint.sh.

servers Directory: Includes a service template with directories and files for building services (controller, db, middleware, schema, server, utils) and a Dockerfile and main.go file for service execution.

## How to Start

### Clone the Repo

```bash
git clone https://github.com/GiveGetGo/givegetgo-backend.git
```

### Create all the env files in each service and redis

### run docker

```bash
docker compose up -d --build
```
