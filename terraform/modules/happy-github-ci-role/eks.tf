//eks:DescribeCluster on resource: arn:aws:eks:us-west-2:401986845158:cluster/si-playground-eks-v2
/*{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": "eks:DescribeCluster",
            "Resource": "arn:aws:eks:us-west-2:401986845158:cluster/si-playground-eks-v2"
        }
    ]
}*/

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
resource "aws_iam_role_policy" "eks" {
  name   = "gh_actions_eks_describe_cluster_${local.namespace}"
  policy = data.aws_iam_policy_document.eks.json
  role   = local.role_name

  depends_on = [module.gh_actions_role]
}

//[FATAL]: unable to initialize the happy client: unable to retrieve integration secret: Unauthorized
// PATCH ADD the following to aws-config-map
/*- "groups":
  - "system:masters"
  "rolearn": "arn:aws:iam::401986845158:role/gh_actions_si_rdev_happy_eks_rdev_happy"
  "username": "gh_actions_si_rdev_happy_eks_rdev_happy"
  */
// TODO still not working exactly
data "kubernetes_config_map" "aws-auth" {
  metadata {
    name      = "aws-auth"
    namespace = "kube-system"
  }
}

locals {
  aws_auth_roles = [{
    groups   = "system:masters"
    rolearn  = module.gh_actions_role.role.arn
    username = module.gh_actions_role.role.name
  }]
  aws_auth_users    = []
  aws_auth_accounts = []
  merged_aws_auth = {
    mapRoles = yamlencode(concat(yamldecode(data.kubernetes_config_map.aws-auth.data.mapRoles), local.aws_auth_roles))
  }
}

resource "kubernetes_config_map_v1_data" "aws-auth" {
  force = true
  metadata {
    name      = "aws-auth"
    namespace = "kube-system"
  }
  data = local.merged_aws_auth

  depends_on = [data.kubernetes_config_map.aws-auth, module.gh_actions_role]
}