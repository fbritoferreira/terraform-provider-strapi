terraform {
  required_providers {
    strapi = {
      source  = "fbritoferreira/strapi"
      version = ">= 0.1.0"
    }
  }
}

provider "strapi" {
  endpoint  = var.strapi_endpoint
  api_token = var.strapi_api_token
}

# Get all available roles in the Strapi instance
data "strapi_roles" "all" {}

# Find role IDs by name for easier reference
locals {
  # Get Super Admin role ID
  super_admin_role_id = [
    for role in data.strapi_roles.all.roles : role.id
    if role.name == "Super Admin"
  ][0]

  # Get Editor role ID (if exists, otherwise use Public)
  editor_role_id = [
    for role in data.strapi_roles.all.roles : role.id
    if role.name == "Editor"
  ][0]

  # Fallback to authenticated role if editor doesn't exist
  default_role_id = [
    for role in data.strapi_roles.all.roles : role.id
    if role.name == "Authenticated"
  ][0]
}

# Create a Super Admin user
resource "strapi_admin_user" "super_admin" {
  email     = "superadmin@example.com"
  firstname = "Super"
  lastname  = "Admin"
  password  = "SuperPass123"

  # Assign Super Admin role
  roles = [tonumber(local.super_admin_role_id)]
}

# Create an Editor user
resource "strapi_admin_user" "editor" {
  email     = "editor@example.com"
  firstname = "Jane"
  lastname  = "Doe"
  password  = "EditorPass123"

  prefered_language = "en"
  is_active         = true

  # Try to assign Editor role, fallback to default
  roles = [
    tonumber(coalesce(local.editor_role_id, local.default_role_id))
  ]
}

# Create another admin user with custom settings
resource "strapi_admin_user" "content_manager" {
  email     = "content@example.com"
  firstname = "Bob"
  lastname  = "Smith"
  password  = "ContentPass123"

  prefered_language = "fr"
  is_active         = true

  # Use default role (Authenticated or similar)
  roles = [tonumber(local.default_role_id)]
}

# Output the created users and available roles
output "admin_users" {
  description = "Created admin users"
  value = {
    super_admin     = strapi_admin_user.super_admin.email
    editor          = strapi_admin_user.editor.email
    content_manager = strapi_admin_user.content_manager.email
  }
}

output "available_roles" {
  description = "All available roles in Strapi"
  value = {
    for role in data.strapi_roles.all.roles : role.name => role.id
  }
}

output "login_credentials" {
  description = "Login credentials for created users"
  sensitive   = true
  value = {
    super_admin = {
      email    = strapi_admin_user.super_admin.email
      password = "SuperPass123"
    }
    editor = {
      email    = strapi_admin_user.editor.email
      password = "EditorPass123"
    }
    content_manager = {
      email    = strapi_admin_user.content_manager.email
      password = "ContentPass123"
    }
  }
}
