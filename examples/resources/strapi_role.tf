terraform {
  required_providers {
    strapi = {
      source  = "fbritoferreira/strapi"
      version = "0.1.0"
    }
  }
}

provider "strapi" {
  endpoint  = "http://localhost:1337"
  api_token = "your-api-token-here"
}

# Create multiple roles
resource "strapi_role" "content_manager" {
  name        = "Content Manager"
  description = "Manages all content types"
}

resource "strapi_role" "author" {
  name        = "Author"
  description = "Can create and edit own content"
}

# Create users with different roles
resource "strapi_user" "admin_user" {
  username  = "admin_user"
  email     = "admin@example.com"
  password  = "AdminSecurePassword123!"
  confirmed = true
  blocked   = false
  role_name = "Authenticated"
}

resource "strapi_user" "content_manager" {
  username  = "content_manager_user"
  email     = "manager@example.com"
  password  = "ManagerSecurePassword123!"
  confirmed = true
  blocked   = false
  role_id   = strapi_role.content_manager.id
}

resource "strapi_user" "author_user" {
  username  = "author_user"
  email     = "author@example.com"
  password  = "AuthorSecurePassword123!"
  confirmed = true
  blocked   = false
  role_id   = strapi_role.author.id
}

# Create a user with role ID (instead of role_name)
resource "strapi_user" "viewer_user" {
  username  = "viewer_user"
  email     = "viewer@example.com"
  password  = "ViewerSecurePassword123!"
  confirmed = true
  blocked   = false
  role_name = "Public"
}
