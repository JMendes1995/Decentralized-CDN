variable "tfstate_bucket_name" {}
variable "project_name" {}
variable "project_id" {}
variable "region" {}
variable "path_local_public_key" {
  sensitive = true
}
variable "username" {}
variable "scopes" {}
variable "image" {}

variable "ip_isp_pub" {

}
variable "webserver_machine_type" {}
variable "webserver_provisioning_model" {}
variable "webserver_tags" {}
variable "webapp_machine_type" {}
variable "webapp_provisioning_model" {}
variable "webapp_tags" {}
variable "db_machine_type" {}
variable "db_provisioning_model" {}
variable "db_tags" {}
variable "rproxy_machine_type" {}
variable "rproxy_provisioning_model" {}
variable "rproxy_tags" {}
