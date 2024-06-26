# Pin specific version for stability
# Use separate stage for building image
# Use debian for easier build utilities
FROM golang:1.22.4-bullseye AS build

# Add non root user
RUN useradd -u 1001 nonroot

WORKDIR /app 

# Copy only files required to install dependencies (better layer caching)
COPY go.mod go.sum ./

# Use cache mount to speed up install of existing dependencies
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

COPY . .

# Compile application during build rather than at runtime
# Add flags to statically link binary
RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o app-golang \
  ./cmd/web

# Use separate stage for deployable image
FROM scratch

WORKDIR /

# Copy the passwd file
COPY --from=build /etc/passwd /etc/passwd

# Copy the binary from the build stage
COPY --from=build /app/app-golang app-golang

COPY --from=build /app/tls /tls

# Use nonroot user
USER nonroot

# Indicate expected port
EXPOSE 8080

CMD ["/app-golang"]
