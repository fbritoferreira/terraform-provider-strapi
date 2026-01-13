# Strapi Terraform Provider Sandbox

This sandbox environment allows you to test the Terraform Strapi provider with a real Strapi instance.

## Prerequisites

- Docker and Docker Compose
- Terraform >= 1.0
- Go (to build the provider locally)

## Quick Start

### 1. Build the Provider

```bash
# From the root directory
make dev
```

This builds and installs the provider to your local Terraform plugin directory.

### 2. Start Strapi

```bash
cd sandbox
docker-compose up -d
```

Wait for Strapi to be fully started (check `docker-compose logs -f`). The admin panel will be available at http://localhost:1337/admin.

### 3. Initialize Strapi

Open your browser to http://localhost:1337/admin and create the first admin user:
- Email: `admin@example.com`
- First name: `Admin`
- Last name: `User`
- Password: `Admin123!`

### 4. Create an API Token

After logging in:
1. Go to **Settings** > **API Tokens** (left sidebar)
2. Click **Create new API Token**
3. Name it: `terraform`
4. Duration: `Unlimited`
5. Token type: `Full access`
6. Click **Save**
7. **Copy the token** - you'll need it for Terraform

### 5. Configure Terraform

```bash
cd sandbox

# Create terraform.tfvars file with your API token
cat > terraform.tfvars <<EOF
strapi_api_token = "your-copied-api-token-here"
strapi_endpoint   = "http://localhost:1337"
EOF
```

### 6. Initialize and Apply

```bash
terraform init
terraform plan
terraform apply
```

## What This Example Creates

The `main.tf` file in this directory demonstrates:

1. **Reading existing roles** - Uses `strapi_roles` data source to get all available roles
2. **Creating admin users** - Creates new admin users and assigns them to specific roles
3. **Role-based assignment** - Finds role IDs by name and assigns to users

## Files

- `docker-compose.yml` - Strapi container configuration
- `main.tf` - Terraform example configuration
- `terraform.tfvars` - Variables file (you need to create this)
- `README.md` - This file

## Testing Specific Scenarios

### Create a Super Admin User

The example creates a super admin user that can manage everything:
```hcl
resource "strapi_admin_user" "super_admin" {
  email    = "superadmin@example.com"
  firstname = "Super"
  lastname  = "Admin"
  password  = "SuperPass123"

  roles = [tonumber(local.super_admin_role_id)]
}
```

### Create an Editor User

The example also creates an editor with limited permissions:
```hcl
resource "strapi_admin_user" "editor" {
  email    = "editor@example.com"
  firstname = "John"
  lastname  = "Editor"
  password  = "EditorPass123"

  roles = [tonumber(local.editor_role_id)]
}
```

## Troubleshooting

### Provider Not Found

If Terraform can't find the provider:
```bash
# Reinstall the provider
cd ..
make dev
cd sandbox
terraform init
```

### Strapi Not Responding

Check if the container is running:
```bash
docker-compose ps
docker-compose logs strapi
```

### API Token Issues

- Ensure the token has **Full access** permissions
- Make sure the token hasn't expired
- Verify the token is correctly copied without extra spaces

### Role Not Found

If you get an error about missing roles:
```bash
# Apply without resources to just see available roles
terraform apply -target=data.strapi_roles.all
```

Then check the output and adjust role names in your configuration.

### Reset the Sandbox

To start fresh:
```bash
docker-compose down -v
docker-compose up -d
terraform destroy
```

## Accessing Strapi

After applying the Terraform configuration, you can log in to the Strapi admin panel at http://localhost:1337/admin with any of the users created by Terraform.

Example credentials (from the example):
- `admin@example.com` / `SuperPass123`
- `editor@example.com` / `EditorPass123`

## Cleanup

```bash
# Destroy Terraform resources
terraform destroy

# Stop Strapi
docker-compose down -v

# Clean up provider (optional)
cd ..
make clean
```

## Development

To make changes to the provider and test them:

1. Edit provider code in `internal/provider/`
2. Rebuild: `make dev`
3. In sandbox: `terraform apply -refresh-only`
4. Apply changes: `terraform apply`

This rebuilds and installs the provider without requiring manual steps.
