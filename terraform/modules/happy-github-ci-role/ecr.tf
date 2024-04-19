data "aws_caller_identity" "current" {}

locals {
  account_id = var.aws_account_id == "" ? data.aws_caller_identity.current.account_id : var.aws_account_id
}

resource "random_pet" "this" {
  keepers = {
    role_name = var.gh_actions_role_name
  }
}

module "ecr_writer_policy" {
  count               = length(var.ecrs) > 0 ? 1 : 0
  source              = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/aws-iam-policy-ecr-writer?ref=v0.125.0"
  role_name           = var.gh_actions_role_name
  ecr_repository_arns = flatten([for ecr in var.ecrs : ecr.repository_arn])
  policy_name         = "gh_actions_ecr_push_${random_pet.this.id}"

  project = var.tags.project
  env     = var.tags.env
  service = var.tags.service
  owner   = var.tags.owner
}

// used for the dynamic autocreated ECRs
module "autocreated_ecr_writer_policy" {
  source    = "git@github.com:chanzuckerberg/shared-infra//terraform/modules/aws-iam-policy-ecr-writer?ref=v0.125.0"
  role_name = var.gh_actions_role_name
  // TODO: not a super fan of this. Would be ideal to have the role only have access to the stacks created by this happy project
  ecr_repository_arns = ["arn:aws:ecr:us-west-2:${local.account_id}:repository/*/${var.tags.env}/*"]
  policy_name         = "gh_actions_ecr_push_${random_pet.this.id}"

  project = var.tags.project
  env     = var.tags.env
  service = var.tags.service
  owner   = var.tags.owner
}

data "aws_iam_policy_document" "ecr_scanner" {
  statement {
    sid = "ScanECR"

    actions = [
      "ecr:BatchGetRepositoryScanningConfiguration",
      "ecr:GetRegistryScanningConfiguration",
      "ecr:DescribeImageScanFindings",
      "ecr:StartImageScan",
      "ecr:PutImageScanningConfiguration",
      "ecr:PutRegistryScanningConfiguration",
      "ecr:PutImageTagMutability"
    ]

    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "ecr_scanner" {
  role        = var.gh_actions_role_name
  name_prefix = "gh_actions_ecr_scan_${random_pet.this.id}"
  policy      = data.aws_iam_policy_document.ecr_scanner.json
}

data "aws_iam_policy_document" "pull_through_cache" {
  statement {
    sid = "PullThroughCacheCorePlatformProdECR"

    actions = [
      "ecr:BatchCheckLayerAvailability",
      "ecr:BatchGetImage",
      "ecr:BatchImportUpstreamImage",
      "ecr:CreateRepository",
      "ecr:DescribeImageScanFindings",
      "ecr:DescribeImages",
      "ecr:DescribeRepositories",
      "ecr:GetAuthorizationToken",
      "ecr:GetDownloadUrlForLayer",
      "ecr:GetLifecyclePolicy",
      "ecr:GetLifecyclePolicyPreview",
      "ecr:GetRepositoryPolicy",
      "ecr:ListImages",
      "ecr:ListTagsForResource",
      "ecr:TagResource",
    ]

    resources = ["arn:aws:ecr:us-west-2:533267185808:repository/*"]
  }
}

resource "aws_iam_role_policy" "pull_through_cache" {
  role        = var.gh_actions_role_name
  name_prefix = "read_only_pull_through_cache_core_platform_prod_access"
  policy      = data.aws_iam_policy_document.pull_through_cache.json
}