variable "gcp_project_id" {
  type        = string
  description = "GCP project id"
}

variable "gcp_region" {
  type        = string
  description = "GCP project region"
}

variable "gcr_image_name" {
  type        = string
  description = "GCR image name"
}

variable "cloud_run_name" {
  type        = string
  description = "Google cloud run name"
}

variable "cloud_run_location" {
  type        = string
  description = "Google cloud run location"
}

variable "ir_cache_dialect" {
  type        = string
  description = "Cache dialect"
}

variable "ir_cache_address" {
  type        = string
  description = "Cache address"
}

variable "ir_cache_password" {
  type        = string
  description = "Cache password"
}

variable "ir_cache_time" {
  type        = string
  description = "Cache time"
}

variable "ir_db_address" {
  type        = string
  description = "Database address"
}

variable "ir_db_name" {
  type        = string
  description = "Database name"
}

variable "ir_db_user" {
  type        = string
  description = "Database user"
}

variable "ir_db_password" {
  type        = string
  description = "Database password"
}

variable "ir_jwt_access_secret" {
  type        = string
  description = "JWT access secret"
}

variable "ir_jwt_access_expired" {
  type        = string
  description = "JWT access expired"
}

variable "ir_jwt_refresh_secret" {
  type        = string
  description = "JWT refresh secret"
}

variable "ir_jwt_refresh_expired" {
  type        = string
  description = "JWT refresh expired"
}

variable "ir_log_json" {
  type        = bool
  description = "Log json"
}

variable "ir_log_level" {
  type        = number
  description = "Log level"
}

variable "ir_newrelic_license_key" {
  type        = string
  description = "Newrelic license key"
}
