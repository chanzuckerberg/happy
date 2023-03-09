locals {
  k8s_namespace = "${var.tags.project}-${var.tags.env}-${var.tags.service}-happy-env"
}

resource "kubernetes_namespace" "happy" {
  metadata {
    name = local.k8s_namespace
  }
}
