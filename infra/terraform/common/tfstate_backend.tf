terraform {
  backend "gcs" {
    bucket  = "cdn-tfstate"
    encrypt = true
  }
}

