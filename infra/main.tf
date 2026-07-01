# ── IAM ─────────────────────────────────────────────────────────────────────

data "aws_iam_policy_document" "ecs_assume" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "execution" {
  name               = "${local.name_prefix}-exec"
  assume_role_policy = data.aws_iam_policy_document.ecs_assume.json
}

resource "aws_iam_role_policy_attachment" "execution_managed" {
  role       = aws_iam_role.execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

data "aws_iam_policy_document" "secrets_read" {
  statement {
    effect    = "Allow"
    actions   = ["ssm:GetParameters", "ssm:GetParameter"]
    resources = [for name in local.secret_ssm_names : "arn:aws:ssm:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:parameter${name}"]
  }
}

resource "aws_iam_role_policy" "execution_secrets" {
  name   = "${local.name_prefix}-secrets"
  role   = aws_iam_role.execution.id
  policy = data.aws_iam_policy_document.secrets_read.json
}

resource "aws_iam_role" "task" {
  name               = "${local.name_prefix}-task"
  assume_role_policy = data.aws_iam_policy_document.ecs_assume.json
}

# ── ECR ────────────────────────────────────────────────────────────────────────

module "ecr" {
  source           = "./modules/ecr"
  repository_names = [for c in local.components : "${var.project}-backend-${c}"]
}

# ── Services ───────────────────────────────────────────────────────────────────

module "api" {
  source = "./modules/service"

  name                   = "${local.name_prefix}-api"
  cluster_arn            = local.cluster_arn
  capacity_provider_name = local.capacity_provider_name
  subnet_ids             = local.subnet_ids
  security_group_id      = local.security_group_id
  execution_role_arn     = aws_iam_role.execution.arn
  task_role_arn          = aws_iam_role.task.arn

  image          = "${local.ecr_registry}/${var.project}-backend-api:${var.image_tag}"
  cpu            = var.api_cpu
  memory         = var.api_memory
  desired_count  = var.api_desired_count
  container_port = var.api_port
  aws_region     = data.aws_region.current.name

  service_connect_namespace = local.service_connect_namespace

  environment_variables = {
    APP_ENV  = var.environment
    PORT     = tostring(var.api_port)
    LOG_JSON = "true"
  }

  secrets = {
    DATABASE_URL            = local.secret_arns["database_url"]
    TIMESERIES_DATABASE_URL = local.secret_arns["timeseries_database_url"]
  }
}

module "worker" {
  source = "./modules/service"

  name                   = "${local.name_prefix}-worker"
  cluster_arn            = local.cluster_arn
  capacity_provider_name = local.capacity_provider_name
  subnet_ids             = local.subnet_ids
  security_group_id      = local.security_group_id
  execution_role_arn     = aws_iam_role.execution.arn
  task_role_arn          = aws_iam_role.task.arn

  image         = "${local.ecr_registry}/${var.project}-backend-worker:${var.image_tag}"
  cpu           = var.worker_cpu
  memory        = var.worker_memory
  desired_count = var.worker_desired_count
  aws_region    = data.aws_region.current.name

  service_connect_namespace = local.service_connect_namespace

  environment_variables = {
    APP_ENV  = var.environment
    LOG_JSON = "true"
  }

  secrets = {
    DATABASE_URL            = local.secret_arns["database_url"]
    TIMESERIES_DATABASE_URL = local.secret_arns["timeseries_database_url"]
    RABBITMQ_URL            = local.secret_arns["rabbitmq_url"]
  }
}

# ── Traefik routing (owned by this app) ──────────────────────────────────────
# The shared Traefik (infra repo) consumes these route declarations. Each app
# publishes its own routes so cluster-wide Traefik settings stay in the infra
# repo while per-app routing lives with the app.
resource "aws_ssm_parameter" "traefik_api_route" {
  name = "/torque/${var.environment}/traefik/routes/backend-api"
  type = "String"
  value = jsonencode({
    host    = var.api_domain
    service = module.api.service_connect_dns
    port    = var.api_port
  })
}
