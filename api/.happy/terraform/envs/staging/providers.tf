provider aws {
  region  = "us-west-2"
  assume_role {
    role_arn = "arn:aws:iam::${var.aws_account_id}:role/${var.aws_role}"
  }
  allowed_account_ids = [var.aws_account_id]
}

provider "aws" {
  alias  = "czi-si"
  region = "us-west-2"

  assume_role {
    role_arn = "arn:aws:iam::626314663667:role/tfe-si"
  }
  allowed_account_ids = ["626314663667"]
}

data "aws_ssm_parameter" "dd_app_key" {
  name     = "/shared-infra-prod-datadog/app_key"
  provider = aws.czi-si
}
data "aws_ssm_parameter" "dd_api_key" {
  name     = "/shared-infra-prod-datadog/api_key"
  provider = aws.czi-si
}

provider "datadog" {
  app_key = data.aws_ssm_parameter.dd_app_key.value
  api_key = data.aws_ssm_parameter.dd_api_key.value
}
