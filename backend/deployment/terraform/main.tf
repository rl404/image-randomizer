terraform {
  backend "gcs" {
    prefix = "state"
  }
}

provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
}
