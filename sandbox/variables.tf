variable "strapi_endpoint" {
  description = "The URL of your Strapi instance"
  type        = string
  default     = "http://localhost:1337"
}

variable "strapi_api_token" {
  description = "Your Strapi API token with Full Access permissions"
  type        = string
  sensitive   = true
}
