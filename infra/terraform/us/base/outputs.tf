output vpc_id {
  value   = module.NetworkUS.vpc_id
}
output private_subnet_name {
  value = local.private_subnet_name
}
