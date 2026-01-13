# Admin User Resource

The `strapi_admin_user` resource manages Strapi admin dashboard users. These are users who can access the Strapi admin panel, as opposed to the content API users managed by `strapi_user`.

## Admin Users API Endpoints

Based on the Strapi source code, the admin users API uses the following endpoints:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/admin/users` | List all admin users (with pagination) |
| POST | `/admin/users` | Create a new admin user |
| GET | `/admin/users/:id` | Get a specific admin user by ID |
| PUT | `/admin/users/:id` | Update an admin user |
| DELETE | `/admin/users/:id` | Delete an admin user |
| POST | `/admin/users/batch-delete` | Delete multiple admin users |

### Authentication

All `/admin/users` endpoints require Bearer token authentication with the admin API token.

### Request/Response Schema

**Create Request (POST /admin/users):**
```json
{
  "email": "user@example.com",
  "firstname": "John",
  "lastname": "Doe",
  "password": "SecurePass123",
  "roles": [1],
  "preferedLanguage": "en"
}
```

**User Response:**
```json
{
  "data": {
    "id": 1,
    "email": "user@example.com",
    "firstname": "John",
    "lastname": "Doe",
    "isActive": true,
    "roles": [1],
    "preferedLanguage": "en",
    "registrationToken": "xxx-xxx-xxx"
  }
}
```

## Schema Reference

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| `email` | string | **Yes** | Email address of the admin user |
| `firstname` | string | **Yes** | First name of the admin user |
| `lastname` | string | No | Last name of the admin user |
| `password` | string | No* | Password for the admin user. Required when creating new users. Must be at least 8 characters with 1 uppercase, 1 lowercase, and 1 digit. |
| `is_active` | bool | No | Whether the admin user account is active. Defaults to `true`. |
| `roles` | list(int) | **Yes** | List of role IDs to assign to the admin user |
| `prefered_language` | string | No | Preferred language for the admin user |

| Attribute | Type | Computed | Description |
|-----------|------|-----------|-------------|
| `id` | string | Yes | The unique ID of the admin user |
| `email` | string | Yes | Email address |
| `firstname` | string | Yes | First name |
| `lastname` | string | Yes | Last name |
| `is_active` | bool | Yes | Whether account is active |
| `roles` | list(int) | Yes | List of role IDs |
| `prefered_language` | string | Yes | Preferred language |
| `registration_token` | string | Yes | Registration token (sensitive) |

## Usage Example

```hcl
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

# Get all available roles in your Strapi instance
data "strapi_roles" "all" {}

# Find role IDs by name
locals {
  super_admin_role_id = [
    for role in data.strapi_roles.all.roles : role.id
    if role.name == "Super Admin"
  ][0]

  editor_role_id = [
    for role in data.strapi_roles.all.roles : role.id
    if role.name == "Editor"
  ][0]
}

# Create an admin user with super admin role
resource "strapi_admin_user" "super_admin" {
  email    = "admin@example.com"
  firstname = "Super"
  lastname  = "Admin"

  # Reference role by ID from data source
  roles = [tonumber(local.super_admin_role_id)]

  password = "SecurePass123"
}

# Create an editor admin user
resource "strapi_admin_user" "editor" {
  email    = "editor@example.com"
  firstname = "Jane"
  lastname  = "Smith"

  # Reference role by ID from data source
  roles = [tonumber(local.editor_role_id)]

  prefered_language = "en"
  is_active          = true
}
```

### Finding Roles

To use role names instead of hardcoded IDs, use the `strapi_roles` data source:

```hcl
data "strapi_roles" "all" {}

locals {
  # Extract role IDs by name for easy reference
  super_admin_role_id = [for role in data.strapi_roles.all.roles : role.id if role.name == "Super Admin"][0]
  editor_role_id    = [for role in data.strapi_roles.all.roles : role.id if role.name == "Editor"][0]
  author_role_id   = [for role in data.strapi_roles.all.roles : role.id if role.name == "Author"][0]
}

resource "strapi_admin_user" "example" {
  email    = "user@example.com"
  firstname = "John"
  lastname  = "Doe"

  roles = [tonumber(local.editor_role_id)]
}
```

## Password Requirements

When creating a new admin user, the password must meet the following criteria (as defined by Strapi's validation):

- Minimum 8 characters
- Maximum 73 bytes
- At least 1 uppercase letter
- At least 1 lowercase letter
- At least 1 digit

## Important Notes

1. **Sensitive Data**: The `password` and `registration_token` fields are marked as sensitive and won't be displayed in plan output.

2. **Password Updates**: You can update the password by providing a new value. After the user is created/updated, the password field is not stored in state for security.

3. **Self-Deletion**: Strapi prevents users from deleting their own account.

4. **Role Assignment**: You must provide at least one role ID. These role IDs correspond to the roles defined in your Strapi instance.

5. **Entity Type**: Admin users are stored in the `admin::user` entity in Strapi's internal database, separate from content API users.

6. **Email Uniqueness**: Email addresses must be unique across all admin users.

## Import

You can import an existing admin user:

```hcl
terraform import strapi_admin_user.example 1
```

Where `1` is the ID of the admin user in Strapi.
