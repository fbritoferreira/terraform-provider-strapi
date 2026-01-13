#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🛑 Stopping Strapi Terraform Provider Sandbox"
echo ""

# Stop Docker containers
echo "🐳 Stopping Docker containers..."
docker-compose down -v
echo "✅ Containers stopped and volumes removed"
echo ""

# Check if terraform.tfstate exists
if [ -f "$SCRIPT_DIR/terraform.tfstate" ]; then
    echo "🗑️  Destroying Terraform resources..."
    terraform destroy -auto-approve 2>/dev/null || true
    echo "✅ Terraform resources destroyed"
    echo ""

    # Clean up Terraform files
    echo "🧹 Cleaning up Terraform files..."
    rm -f "$SCRIPT_DIR/terraform.tfstate"*
    rm -f "$SCRIPT_DIR/.terraform.lock.hcl"
    rm -rf "$SCRIPT_DIR/.terraform/"
    echo "✅ Terraform files cleaned"
fi

echo ""
echo "✨ Sandbox reset complete!"
echo "To start again: ./start.sh"
