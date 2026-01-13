terraform {
  required_providers {
    strapi = {
      source  = "fbritoferreira/strapi"
      version = ">= 0.1.0"
    }
  }
}

provider "strapi" {
  endpoint = "http://localhost:1337"
  api_token = var.strapi_api_token
}

variable "strapi_api_token" {
  description = "Your Strapi API token"
  type        = string
  sensitive   = true
}

# Get all available roles to reference by name
data "strapi_roles" "all" {}

# Find super admin role ID
locals {
  super_admin_role_id = [for role in data.strapi_roles.all.roles : role.id if role.name == "Super Admin"][0]
  editor_role_id = [for role in data.strapi_roles.all.roles : role.id if role.name == "Editor"][0]
}

# Create an admin user with super admin role
resource "strapi_admin_user" "admin" {
  email    = "admin@example.com"
  firstname = "Admin"
  lastname  = "User"

  roles = [tonumber(local.super_admin_role_id)]

  password = "SecurePass123"
}

# Create another admin user with editor role
resource "strapi_admin_user" "editor" {
  email    = "editor@example.com"
  firstname = "John"
  lastname  = "Doe"

  roles = [tonumber(local.editor_role_id)]

  prefered_language = "en"
  is_active          = true
}
  }
}

provider "strapi" {
  endpoint  = "http://localhost:1337"
  api_token = "your-api-token-here"
}

# Create an admin user
resource "strapi_admin_user" "admin" {
  email     = "admin@example.com"
  firstname = "Admin"
  lastname  = "User"

  roles = [1]

  password = "SecurePass123"
}

# Create another admin user with custom settings
resource "strapi_admin_user" "editor" {
  email     = "editor@example.com"
  firstname = "John"
  lastname  = "Doe"

  roles = [2]

  prefered_language = "en"
  is_active         = true
}
