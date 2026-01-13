# Terraform Provider for Strapi

A Terraform provider for managing Strapi CMS resources.

## Features

- [x] Manage content API users
- [x] Manage admin dashboard users
- [x] Manage roles
- [x] Query available roles as data source
- [ ] Manage content types
- [ ] Manage collection types
- [ ] Manage API tokens
- [ ] Manage media library

## Installation

### From Terraform Registry

```hcl
terraform {
  required_providers {
    strapi = {
      source  = "fbritoferreira/strapi"
      version = "~> 0.1"
    }
  }
}
```

### Local Development

See [DEVELOPMENT.md](DEVELOPMENT.md)

### Testing with Sandbox

A complete sandbox environment is provided in the `sandbox/` directory. It includes:

- Docker Compose configuration for running a Strapi instance
- Example Terraform configurations
- Scripts to quickly set up and tear down the test environment

To get started with the sandbox:

```bash
cd sandbox
./start.sh
```

See [sandbox/README.md](sandbox/README.md) for detailed instructions.

## Configuration

```hcl
provider "strapi" {
  endpoint  = "http://localhost:1337"
  api_token = var.strapi_api_token
}
```

### Environment Variables

- `STRAPI_ENDPOINT` - Strapi API endpoint URL
- `STRAPI_API_TOKEN` - Strapi API token for authentication

## Example Usage

### Managing Admin Users

```hcl
terraform {
  required_providers {
    strapi = {
      source  = "fbritoferreira/strapi"
      version = "~> 0.1"
    }
  }
}

provider "strapi" {
  endpoint  = "http://localhost:1337"
  api_token = var.strapi_api_token
}

# Get all available roles
data "strapi_roles" "all" {}

# Find Super Admin role ID
locals {
  super_admin_role_id = [
    for role in data.strapi_roles.all.roles : role.id
    if role.name == "Super Admin"
  ][0]
}

# Create a super admin user
resource "strapi_admin_user" "admin" {
  email    = "admin@example.com"
  firstname = "Admin"
  lastname  = "User"
  password  = "SecurePass123"

  roles = [tonumber(local.super_admin_role_id)]
}
```

### Managing Content API Users and Roles

```hcl
terraform {
  required_providers {
    strapi = {
      source  = "fbritoferreira/strapi"
      version = "~> 0.1"
    }
  }
}

provider "strapi" {
  endpoint  = "http://localhost:1337"
  api_token = var.strapi_api_token
}

# Create a custom role
resource "strapi_role" "editor" {
  name        = "Editor"
  description = "Users who can edit content"
}

# Create a user and assign to role
resource "strapi_user" "john_doe" {
  username  = "john_doe"
  email     = "john.doe@example.com"
  password  = "SecurePassword123!"
  confirmed  = true
  blocked    = false
  role_name  = strapi_role.editor.name
}

# Or use role_id instead of role_name
resource "strapi_user" "jane_doe" {
  username  = "jane_doe"
  email     = "jane.doe@example.com"
  password  = "SecurePassword123!"
  confirmed  = true
  blocked    = false
  role_id   = strapi_role.editor.id
}
```

### Available Resources

- **strapi_user**: Manage Strapi content API users
- **strapi_admin_user**: Manage Strapi admin dashboard users
- **strapi_role**: Manage Strapi roles

### Available Data Sources

- **strapi_roles**: Query all available roles in Strapi

## Contributing

Contributions are welcome! Please read [DEVELOPMENT.md](DEVELOPMENT.md) for guidelines.

## License

Apache 2.0
