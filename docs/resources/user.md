# `strapi_user`

Manages a Strapi user.

## Example Usage

```hcl
resource "strapi_user" "example" {
  username  = "john_doe"
  email     = "john.doe@example.com"
  password  = "SecurePassword123!"
  confirmed  = true
  blocked    = false
  role_name  = "Authenticated"
}
```

## Argument Reference

The following arguments are supported:

- `username` - (Required) The username of the user.
- `email` - (Required) The email address of the user.
- `password` - (Optional, Sensitive) The password for the user. Only used when creating a new user.
- `confirmed` - (Optional) Whether the user account is confirmed. Defaults to `false`.
- `blocked` - (Optional) Whether the user account is blocked. Defaults to `false`.
- `role_name` - (Optional) The name of the role to assign to the user. Mutually exclusive with `role_id`.
- `role_id` - (Optional) The ID of the role assigned to the user. Mutually exclusive with `role_name`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the user.
- `document_id` - The document ID of the user.
- `created_at` - Timestamp when the user was created.
- `updated_at` - Timestamp when the user was last updated.

## Import

Users can be imported using the user ID:

```hcl
terraform import strapi_user.example 1
```
