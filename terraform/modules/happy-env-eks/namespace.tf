resource "kubernetes_namespace" "happy" {
  metadata {
    name = "${var.tags.project}-${var.tags.env}-${var.tags.service}-happy-env"
  }
}
