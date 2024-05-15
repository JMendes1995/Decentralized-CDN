
module "ContentStorageBucket" {
  source        = "../modules/gcp/storage_bucket"
  name          = "cdn-content-storage"
  force_destroy = true
  location      = "EU"
  storage_class = "STANDARD"
  versioning    = true
}


resource "google_storage_default_object_access_control" "public_rule" {
  bucket     = "cdn-content-storage"
  role       = "READER"
  entity     = "allUsers"
  depends_on = [module.ContentStorageBucket]
}


module "DB" {
  source = "../modules/gcp/compute/private_vm"

  num_instances      = 1
  vm_name            = "db"
  machine_type       = var.db_machine_type
  vpc_id             = data.terraform_remote_state.base_tfstate.outputs.vpc_id
  subnet             = data.terraform_remote_state.base_tfstate.outputs.private_subnet_name
  public_instance    = true
  image              = var.image
  provisioning_model = var.db_provisioning_model
  tags               = var.db_tags
  scopes             = var.scopes
  ssh_pub            = file(var.path_local_public_key)
  username           = var.username
  defaul_sa_name     = data.google_compute_default_service_account.default_sa.email
  available_zones    = ["europe-west4-a", "europe-west4-b", "europe-west4-c"]
  packages           = "mariadb-server"
  static_ip          = "10.10.0.2"

}


module "Volumes" {
  source = "../modules/gcp/compute/storage"

  storage_device_number = 1
  storage_device_name   = "db-storage"
  storage_device_type   = "pd-standard"
  storage_device_size   = 5
  available_zones       = ["europe-west4-a", "europe-west4-b", "europe-west4-c"]
}

resource "google_compute_attached_disk" "attached_storage" {
  count = 1
  disk = element(module.Volumes.volume_ids, count.index)
  instance = element(module.DB.vm_ids, count.index)
  zone = element(["europe-west4-a", "europe-west4-b", "europe-west4-c"], count.index)

  lifecycle {
    ignore_changes = [instance]
  }
}
