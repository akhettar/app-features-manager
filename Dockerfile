# Build stage
FROM golang:alpine3.11 AS build

ENV CONFIG_FOLDR /go/config

# Support CGO and SSL
USER root
RUN apk --no-cache add gcc g++ make ca-certificates
USER app

WORKDIR /go/src/app-features-api

# Copy the source code and config
COPY . .
COPY config $CONFIG_FOLDR

# Build binary
RUN go install

# Production build stage
FROM alpine
ENV ENVIRONMENT default
ENV CONFIG_FOLDR /app/config
ENV GIN_MODE release

WORKDIR /app

# Copy built binaries
COPY --from=build /go/bin /app
COPY --from=build /go/config /app/config

HEALTHCHECK --interval=5s --retries=10 CMD curl -fs http://localhost:8080/health || exit 1
EXPOSE 8080

CMD /app/app-features-api
