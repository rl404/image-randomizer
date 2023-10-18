resource "google_cloud_run_v2_service" "server" {
  name     = var.cloud_run_name
  location = var.cloud_run_location
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    labels = {
      app = var.cloud_run_name
    }
    scaling {
      min_instance_count = 0
    }
    containers {
      name    = var.cloud_run_name
      image   = var.gcr_image_name
      command = ["./image-randomizer"]
      args    = ["server"]
      env {
        name  = "IR_CACHE_DIALECT"
        value = var.ir_cache_dialect
      }
      env {
        name  = "IR_CACHE_ADDRESS"
        value = var.ir_cache_address
      }
      env {
        name  = "IR_CACHE_PASSWORD"
        value = var.ir_cache_password
      }
      env {
        name  = "IR_CACHE_TIME"
        value = var.ir_cache_time
      }
      env {
        name  = "IR_DB_ADDRESS"
        value = var.ir_db_address
      }
      env {
        name  = "IR_DB_NAME"
        value = var.ir_db_name
      }
      env {
        name  = "IR_DB_USER"
        value = var.ir_db_user
      }
      env {
        name  = "IR_DB_PASSWORD"
        value = var.ir_db_password
      }
      env {
        name  = "IR_JWT_ACCESS_SECRET"
        value = var.ir_jwt_access_secret
      }
      env {
        name  = "IR_JWT_ACCESS_EXPIRED"
        value = var.ir_jwt_access_expired
      }
      env {
        name  = "IR_JWT_REFRESH_SECRET"
        value = var.ir_jwt_refresh_secret
      }
      env {
        name  = "IR_JWT_REFRESH_EXPIRED"
        value = var.ir_jwt_refresh_expired
      }
      env {
        name  = "IR_LOG_JSON"
        value = var.ir_log_json
      }
      env {
        name  = "IR_LOG_LEVEL"
        value = var.ir_log_level
      }
      env {
        name  = "IR_NEWRELIC_LICENSE_KEY"
        value = var.ir_newrelic_license_key
      }
    }
  }
}

resource "google_cloud_run_service_iam_binding" "noauth" {
  service  = google_cloud_run_v2_service.server.name
  location = google_cloud_run_v2_service.server.location
  role     = "roles/run.invoker"
  members  = ["allUsers"]
}
