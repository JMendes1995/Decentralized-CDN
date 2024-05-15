
module "WebServer" {
  source = "../modules/gcp/compute/private_vm"

  num_instances      = 1
  vm_name            = "webserver"
  machine_type       = var.webserver_machine_type
  vpc_id             = data.terraform_remote_state.base_tfstate.outputs.vpc_id
  subnet             = data.terraform_remote_state.base_tfstate.outputs.private_subnet_name
  public_instance    = true
  image              = var.image
  provisioning_model = var.webserver_provisioning_model
  tags               = var.webserver_tags
  scopes             = var.scopes
  ssh_pub            = file(var.path_local_public_key)
  username           = var.username
  defaul_sa_name     = data.google_compute_default_service_account.default_sa.email
  available_zones    = ["europe-west4-a", "europe-west4-b", "europe-west4-c"]
  packages           = "dnsutils"
  static_ip          = "10.10.0.3"

}
