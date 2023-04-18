data "aws_iam_policy_document" "eks" {
  statement {
    sid = "GhActionsDescribeCluster"
    actions = [
      "eks:DescribeCluster",
    ]
    resources = [
      var.eks_cluster_arn,
    ]
  }
}

resource "random_pet" "this" {
  keepers = {
    role_name = var.gh_actions_role.role.name
  }
}

resource "aws_iam_role_policy" "eks" {
  name   = "gh_actions_eks_describe_cluster_${random_pet.this.id}"
  policy = data.aws_iam_policy_document.eks.json
  role   = random_pet.this.keepers.role_name
}

data "kubernetes_config_map" "aws-auth" {
  metadata {
    name      = "aws-auth"
    namespace = "kube-system"
  }
}

locals {
  aws_auth_roles = [{
    groups   = "system:masters"
    rolearn  = var.gh_actions_role.role.arn
    username = var.gh_actions_role.role.name
  }]
  aws_auth_users    = []
  aws_auth_accounts = []
  merged_aws_auth = {
    mapRoles    = yamlencode(concat(yamldecode(data.kubernetes_config_map.aws-auth.data.mapRoles), local.aws_auth_roles))
    mapUsers    = yamlencode(concat(yamldecode(data.kubernetes_config_map.aws-auth.data.mapUsers), local.aws_auth_users))
    mapAccounts = yamlencode(concat(yamldecode(data.kubernetes_config_map.aws-auth.data.mapAccounts), local.aws_auth_accounts))
  }
}

resource "kubernetes_config_map_v1_data" "aws-auth" {
  force = true
  metadata {
    name      = "aws-auth"
    namespace = "kube-system"
  }
  data = local.merged_aws_auth

  depends_on = [
    data.kubernetes_config_map.aws-auth,
  ]
}