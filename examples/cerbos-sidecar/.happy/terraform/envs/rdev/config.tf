resource "kubernetes_config_map" "cerbos-sidecar-demo" {
  metadata {
    name = "${var.stack_name}-cerbos-sidecar-demo"
    labels = {
      "app.kubernetes.io/name"       = var.stack_name
      "app.kubernetes.io/component"  = var.stack_name
      "app.kubernetes.io/part-of"    = var.stack_name
      "app.kubernetes.io/managed-by" = "happy"
    }
  }
  data = {
    "config.yaml" = <<-EOT
      server:
        # Configure Cerbos to listen on a Unix domain socket.
        httpListenAddr: "unix:/sock/cerbos.sock"
      storage:
        driver: disk
        disk:
          directory: /policies
          watchForChanges: false
    EOT
  }
}

resource "kubernetes_certificate" "cerbos-sidecar-demo" {
  metadata {
    name = "${var.stack_name}-cerbos-sidecar-demo"
    labels = {
      "app.kubernetes.io/name"       = var.stack_name
      "app.kubernetes.io/component"  = var.stack_name
      "app.kubernetes.io/part-of"    = var.stack_name
      "app.kubernetes.io/managed-by" = "happy"
    }
  }
  spec {
    is_ca = true
    secret_name = "cerbos-sidecar-demo"
    dns_names = [
      "cerbos-sidecar-demo.default.svc.cluster.local",
      "cerbos-sidecar-demo.default.svc",
      "cerbos-sidecar-demo.default",
      "cerbos-sidecar-demo"
    ]
    issuer_ref {
      name = "selfsigned-cluster-issuer"
      kind = "ClusterIssuer"
      group = "cert-manager.io"
    }
  }
}