# Use base golang image from Docker Hub
FROM golang:1.17 AS build

WORKDIR /mater-crud-api

# Install dependencies in go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the application source code
COPY . ./

# Compile the application to /app.
# Skaffold passes in debug-oriented compiler flags
ARG SKAFFOLD_GO_GCFLAGS
RUN echo "Go gcflags: ${SKAFFOLD_GO_GCFLAGS}"
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -mod=readonly -v -o /app

# Now create separate deployment image
FROM gcr.io/distroless/base

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
# Default behavior - a failure prints a stack trace for the current goroutine.
# See https://golang.org/pkg/runtime/

ENV GOTRACEBACK=single
# ENV API_SECRET=secret
# ENV DB_HOST=10.54.64.3
# ENV DB_DRIVER=postgres
# ENV DB_NAME=postgres
# ENV DB_USER=postgres
# ENV DB_PASSWORD=passmater
# ENV DB_PORT=5432

# Copy template & assets
WORKDIR /mater-crud-api
COPY --from=build /app ./app

ENTRYPOINT ["./app"]