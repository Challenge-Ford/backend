provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "torque"
      Application = "backend"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }
}
