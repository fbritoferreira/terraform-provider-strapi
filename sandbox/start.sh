#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROVIDER_DIR="$(dirname "$SCRIPT_DIR")"

echo "🚀 Strapi Terraform Provider Sandbox - Quick Start"
echo ""

# Check if terraform is installed
if ! command -v terraform &> /dev/null; then
    echo "❌ Error: terraform is not installed"
    echo "Please install Terraform: https://www.terraform.io/downloads.html"
    exit 1
fi

# Check if docker-compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Error: docker-compose is not installed"
    echo "Please install Docker Compose: https://docs.docker.com/compose/install/"
    exit 1
fi

# Build and install provider
echo "📦 Building and installing Terraform provider..."
cd "$PROVIDER_DIR"
make dev
cd "$SCRIPT_DIR"
echo "✅ Provider installed"
echo ""

# Start Strapi
echo "🐳 Starting Strapi container..."
docker-compose up -d
echo "✅ Strapi container started"
echo ""

# Wait for Strapi to be ready
echo "⏳ Waiting for Strapi to be ready..."
MAX_RETRIES=30
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -s http://localhost:1337/admin > /dev/null 2>&1; then
        echo "✅ Strapi is ready!"
        break
    fi
    echo "   Waiting... ($((RETRY_COUNT + 1))/$MAX_RETRIES)"
    sleep 2
    RETRY_COUNT=$((RETRY_COUNT + 1))
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo "❌ Error: Strapi did not start in time"
    echo "Check logs with: docker-compose logs"
    exit 1
fi

echo ""
echo "📋 Next steps:"
echo "   1. Open http://localhost:1337/admin in your browser"
echo "   2. Create your first admin user"
echo "   3. Go to Settings > API Tokens and create a Full Access token"
echo "   4. Create terraform.tfvars: cp terraform.tfvars.example terraform.tfvars"
echo "   5. Edit terraform.tfvars and add your API token"
echo "   6. Run: terraform init && terraform plan && terraform apply"
echo ""
echo "ℹ️  Strapi Admin: http://localhost:1337/admin"
echo "ℹ️  Strapi API:    http://localhost:1337"
echo ""
echo "To stop Strapi:   docker-compose down"
echo "To reset sandbox: docker-compose down -v && terraform destroy"
