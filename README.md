# golang-aws-secrets-manager-usage-samples
golang-aws-secrets-manager-usage-samples


```
export ORG_ID=byjum ENVIRON=stg PROJECT_NAME=testsecret AWS_DEFAULT_REGION=ap-southeast-1     ## PROJECT_NAME & ORG_ID can be hardcoded
export RUNNING_ON_LOCAL=yes     ## Set this to "no" while deploying the application
go mod tidy
go run .
```

Example - 

```
Looing for secret -  /byjum/stg/testsecret-secrets
```

```
Application running on local. Feed an ENV file using `docker run` or `docker-compose`
```

```
Environment variable RUNNING_ON_LOCAL not set OR it has a wrong value
```