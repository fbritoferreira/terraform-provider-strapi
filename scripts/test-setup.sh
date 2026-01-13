#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
STRAPI_DIR="$PROJECT_DIR/sandbox/strapi-app"
PID_FILE="$PROJECT_DIR/.strapi.pid"
TOKEN_FILE="$PROJECT_DIR/.strapi-test-token"

ADMIN_EMAIL="admin@test.com"
ADMIN_PASSWORD="Admin123!"
ADMIN_FIRSTNAME="Test"
ADMIN_LASTNAME="Admin"

cleanup() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if kill -0 "$PID" 2>/dev/null; then
            echo "Stopping Strapi (PID: $PID)..."
            kill "$PID" 2>/dev/null || true
            sleep 2
            kill -9 "$PID" 2>/dev/null || true
        fi
        rm -f "$PID_FILE"
    fi
    rm -f "$TOKEN_FILE"
}

wait_for_strapi() {
    local max_attempts=60
    local attempt=0
    echo "Waiting for Strapi to be ready..."
    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:1337/_health > /dev/null 2>&1; then
            echo "Strapi is ready!"
            return 0
        fi
        attempt=$((attempt + 1))
        echo "  Attempt $attempt/$max_attempts..."
        sleep 2
    done
    echo "Strapi failed to start"
    return 1
}

create_admin_user() {
    echo "Creating admin user..."
    cd "$STRAPI_DIR"
    
    pnpm strapi admin:create-user \
        --firstname="$ADMIN_FIRSTNAME" \
        --lastname="$ADMIN_LASTNAME" \
        --email="$ADMIN_EMAIL" \
        --password="$ADMIN_PASSWORD" 2>/dev/null || echo "Admin user may already exist"
}

get_jwt_token() {
    echo "Getting JWT token..."
    
    local response=$(curl -s -X POST http://localhost:1337/admin/login \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}")
    
    local token=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    
    if [ -z "$token" ]; then
        echo "Failed to get JWT token. Response: $response"
        return 1
    fi
    
    echo "$token"
}

create_api_token() {
    local jwt_token="$1"
    echo "Creating API token..."
    
    local response=$(curl -s -X POST http://localhost:1337/admin/api-tokens \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $jwt_token" \
        -d '{
            "name": "Test Token",
            "description": "Token for acceptance tests",
            "type": "full-access",
            "lifespan": null
        }')
    
    local api_token=$(echo "$response" | grep -o '"accessKey":"[^"]*"' | cut -d'"' -f4)
    
    if [ -z "$api_token" ]; then
        echo "Failed to create API token. Response: $response"
        return 1
    fi
    
    echo "$api_token"
}

case "${1:-start}" in
    start)
        cleanup
        
        echo "Building Strapi..."
        cd "$STRAPI_DIR"
        pnpm build
        
        echo "Starting Strapi..."
        pnpm start &
        echo $! > "$PID_FILE"
        
        if ! wait_for_strapi; then
            cleanup
            exit 1
        fi
        
        create_admin_user
        
        jwt_token=$(get_jwt_token)
        if [ -z "$jwt_token" ]; then
            echo "Failed to authenticate"
            cleanup
            exit 1
        fi
        
        api_token=$(create_api_token "$jwt_token")
        if [ -z "$api_token" ]; then
            echo "Failed to create API token"
            cleanup
            exit 1
        fi
        
        echo "$api_token" > "$TOKEN_FILE"
        
        echo ""
        echo "Strapi is running!"
        echo "  Endpoint: http://localhost:1337"
        echo "  API Token: $api_token"
        echo "  PID: $(cat "$PID_FILE")"
        echo ""
        echo "To run tests:"
        echo "  export STRAPI_ENDPOINT=http://localhost:1337"
        echo "  export STRAPI_API_TOKEN=$api_token"
        echo "  make testacc"
        ;;
    
    stop)
        cleanup
        echo "Strapi stopped"
        ;;
    
    token)
        if [ -f "$TOKEN_FILE" ]; then
            cat "$TOKEN_FILE"
        else
            echo "No token file found. Run 'start' first."
            exit 1
        fi
        ;;
    
    *)
        echo "Usage: $0 {start|stop|token}"
        exit 1
        ;;
esac
