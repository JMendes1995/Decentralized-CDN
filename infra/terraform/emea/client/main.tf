
module "Client" {
  source = "../modules/gcp/compute/private_vm"

  num_instances      = 2
  vm_name            = "client"
  machine_type       = var.webapp_machine_type
  vpc_id             = data.terraform_remote_state.base_tfstate.outputs.vpc_id
  subnet             = data.terraform_remote_state.base_tfstate.outputs.private_subnet_name
  public_instance    = true
  image              = var.image
  provisioning_model = var.webapp_provisioning_model
  tags               = var.webapp_tags
  scopes             = var.scopes
  ssh_pub            = file(var.path_local_public_key)
  username           = var.username
  defaul_sa_name     = data.google_compute_default_service_account.default_sa.email
  available_zones    = ["europe-west4-a", "europe-west4-b", "europe-west4-c"]
  packages           = "dnsutils memcached libmemcached-tools"

}
