# `strapi_role`

Manages a Strapi role.

## Example Usage

```hcl
resource "strapi_role" "editor" {
  name        = "Editor"
  description = "Users who can edit content"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the role.
- `description` - (Optional) The description of the role.
- `type` - (Optional) The type of the role (e.g., 'authenticated', 'public').

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the role.
- `created_at` - Timestamp when the role was created.
- `updated_at` - Timestamp when the role was last updated.

## Notes

Strapi provides two default roles that cannot be deleted:

- `Public`: for users accessing content without authentication
- `Authenticated`: for logged-in users

You can create additional custom roles for more granular permission control.

## Import

Roles can be imported using the role ID:

```hcl
terraform import strapi_role.editor 3
```
