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
    role_name = var.gh_actions_role_name
  }
}

resource "aws_iam_role_policy" "eks" {
  name   = "gh_actions_eks_describe_cluster_${random_pet.this.id}"
  policy = data.aws_iam_policy_document.eks.json
  role   = random_pet.this.keepers.role_name
}