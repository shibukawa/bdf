# The Print Manager as a function of AWS Lambda (a container image), as
# pw build --target=lambda packages a Popcorn Web application: the AWS
# Lambda Web Adapter turns the invocations into HTTP requests to the
# process. Build it from the repository's root:
#
#   docker build --platform linux/arm64 -f cmd/printmanager/deploy/lambda/Dockerfile -t printmanager .
#
# Configure it with environment variables (no config file is read):
#
#   MANAGER_SERVERLESS=true
#   MANAGER_RECORDS=dynamodb://<table>?region=ap-northeast-1   (keys pk and sk, strings)
#   MANAGER_DOCUMENTS=s3://<bucket>?region=ap-northeast-1      (not the bucket of the archive)
#   MANAGER_SECRET, MANAGER_ADMIN_TOKEN, MANAGER_TICK_TOKEN, MANAGER_SUBMIT_TOKENS
#
# and let EventBridge Scheduler call POST /api/tick (bearer MANAGER_TICK_TOKEN)
# every minute. The function needs a VPC (with endpoints for S3 and
# DynamoDB) to reach printers over Direct Connect or a VPN.
FROM golang:1.27-bookworm AS build
ARG TARGETARCH
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /bootstrap ./cmd/printmanager

FROM public.ecr.aws/awsguru/aws-lambda-adapter:1.0.1 AS adapter

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=adapter /lambda-adapter /opt/extensions/lambda-adapter
COPY --from=build /bootstrap /bootstrap
ENV APP_ENV=prod PORT=8080 MANAGER_SERVERLESS=true TMPDIR=/tmp
USER nonroot
ENTRYPOINT ["/bootstrap"]
