environment      = "staging"
shared_state_key = "staging/terraform.tfstate"

api_domain = "api.staging.torque-next.space"

api_desired_count = 1
api_cpu           = 256
api_memory        = 512

worker_desired_count = 1
worker_cpu           = 256
worker_memory        = 512

# SSM Parameter Store names (create these in the staging account).
database_url_ssm            = "/torque/staging/backend/database_url"
timeseries_database_url_ssm = "/torque/staging/backend/timeseries_database_url"
rabbitmq_url_ssm            = "/torque/staging/backend/rabbitmq_url"
