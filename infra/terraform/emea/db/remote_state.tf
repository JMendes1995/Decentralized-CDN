data "terraform_remote_state" "base_tfstate" {
  backend = "gcs"

  config = {
    bucket  = var.tfstate_bucket_name
    prefix  = "terraform/state/baseemea"
  }
}
