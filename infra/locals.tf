data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

locals {
  name_prefix  = "${var.project}-backend-${var.environment}"
  ecr_registry = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${data.aws_region.current.name}.amazonaws.com"

  components = ["api", "worker"]

  # SSM parameter names holding application secrets.
  secret_ssm_names = {
    database_url            = var.database_url_ssm
    timeseries_database_url = var.timeseries_database_url_ssm
    rabbitmq_url            = var.rabbitmq_url_ssm
  }

  # Full ARNs of those parameters, used by the container `secrets` block.
  secret_arns = {
    for k, name in local.secret_ssm_names :
    k => "arn:aws:ssm:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:parameter${name}"
  }
}
