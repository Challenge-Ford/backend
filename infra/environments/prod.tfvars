environment      = "prod"
shared_state_key = "prod/terraform.tfstate"

api_domain = "api.torque-next.space"

api_desired_count = 2
api_cpu           = 512
api_memory        = 1024

worker_desired_count = 1
worker_cpu           = 512
worker_memory        = 1024

# SSM Parameter Store names (create these in the prod account).
database_url_ssm            = "/torque/prod/backend/database_url"
timeseries_database_url_ssm = "/torque/prod/backend/timeseries_database_url"
rabbitmq_url_ssm            = "/torque/prod/backend/rabbitmq_url"
