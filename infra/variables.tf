variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name for this app stack (staging or prod)"
  type        = string

  validation {
    condition     = contains(["staging", "prod"], var.environment)
    error_message = "environment must be either 'staging' or 'prod'."
  }
}

variable "project" {
  description = "Project prefix, shared with the infra repo"
  type        = string
  default     = "torque"
}

# ── Shared infra (infra repo) reference ───────────────────────────────────────

variable "tf_state_bucket" {
  description = "S3 bucket that holds the shared infra Terraform state"
  type        = string
  default     = "torque-tf-state"
}

variable "tf_locks_table" {
  description = "DynamoDB table used for Terraform state locking"
  type        = string
  default     = "torque-terraform-locks"
}

variable "shared_state_key" {
  description = "State key of the shared infra environment to consume (e.g. staging/terraform.tfstate)"
  type        = string
}

# ── Image ─────────────────────────────────────────────────────────────────────

variable "image_tag" {
  description = "Container image tag to deploy (usually the git SHA)"
  type        = string
}

# ── API service ────────────────────────────────────────────────────────────────

variable "api_domain" {
  description = "Public hostname the API is exposed on through Traefik"
  type        = string
}

variable "api_port" {
  description = "Port the API container listens on"
  type        = number
  default     = 8080
}

variable "api_desired_count" {
  description = "Number of API tasks to run"
  type        = number
  default     = 1
}

variable "api_cpu" {
  description = "API task CPU units"
  type        = number
  default     = 256
}

variable "api_memory" {
  description = "API task memory (MiB)"
  type        = number
  default     = 512
}

# ── Worker service ──────────────────────────────────────────────────────────────

variable "worker_desired_count" {
  description = "Number of worker tasks to run"
  type        = number
  default     = 1
}

variable "worker_cpu" {
  description = "Worker task CPU units"
  type        = number
  default     = 256
}

variable "worker_memory" {
  description = "Worker task memory (MiB)"
  type        = number
  default     = 512
}

# ── Application secrets (SSM parameter names in the target account) ──────────────

variable "database_url_ssm" {
  description = "SSM Parameter Store name holding the Postgres connection string"
  type        = string
}

variable "timeseries_database_url_ssm" {
  description = "SSM Parameter Store name holding the TimescaleDB connection string"
  type        = string
}

variable "rabbitmq_url_ssm" {
  description = "SSM Parameter Store name holding the RabbitMQ connection string"
  type        = string
}
