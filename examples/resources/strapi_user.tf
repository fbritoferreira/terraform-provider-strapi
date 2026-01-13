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

# Create a new role
resource "strapi_role" "editor" {
  name        = "Editor"
  description = "Users who can edit content"
}

# Create a user and assign the editor role
resource "strapi_user" "john_doe" {
  username  = "john_doe"
  email     = "john.doe@example.com"
  password  = "SecurePassword123!"
  confirmed = true
  blocked   = false
  role_name = strapi_role.editor.name
}
