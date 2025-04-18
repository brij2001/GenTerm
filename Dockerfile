ARG LLM_API_KEY
ARG LLM_MODEL
ARG LLM_BASE_URL

FROM node:18 AS frontend-build

WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install

# Copy only the necessary source files, excluding node_modules
COPY frontend/src ./src
COPY frontend/public ./public
COPY frontend/.env* ./
COPY frontend/tsconfig*.json ./
COPY frontend/*.config.js ./
ENV REACT_APP_API_URL=http://localhost:8080
RUN npm run build

FROM golang:1.21-alpine AS backend-build

WORKDIR /app/backend
COPY backend/ ./
RUN go mod download
# Set environment variables for build time
ARG LLM_BASE_URL
ARG LLM_API_KEY
ARG LLM_MODEL
ENV PORT=8080
ENV LLM_BASE_URL=${LLM_BASE_URL}
ENV LLM_API_KEY=${LLM_API_KEY}
ENV LLM_MODEL=${LLM_MODEL}
# Use ldflags to embed these values directly into the binary
RUN go build -ldflags="-X 'main.LlmApiKey=${LLM_API_KEY}' -X 'main.LlmModel=${LLM_MODEL}' -X 'main.LlmBaseUrl=${LLM_BASE_URL}'" -o server ./cmd/server/main.go

FROM alpine:latest

# Install Node.js for serving the frontend
RUN apk add --no-cache nodejs npm

# Create app directory
WORKDIR /app

# Copy frontend build from frontend-build stage
COPY --from=frontend-build /app/frontend/build /app/frontend/build

# Copy the backend executable from backend-build stage
COPY --from=backend-build /app/backend/server /app/server

# Expose port for Azure Web App
EXPOSE 3000

# Create a startup script
RUN echo '#!/bin/sh' > /app/startup.sh && \
    echo 'npm install -g serve' >> /app/startup.sh && \
    echo 'serve -s /app/frontend/build -l 3000 &' >> /app/startup.sh && \
    echo 'cd /app && ./server' >> /app/startup.sh && \
    chmod +x /app/startup.sh

# Set the startup command
CMD ["/app/startup.sh"] 