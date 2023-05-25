resource "kubernetes_manifest" "traefik-ingress-route" {
  manifest = {
    "apiVersion" = "traefik.io/v1alpha1"
    "kind"       = "IngressRoute"
    "metadata" = {
      "name"      = var.ingress_name
      "namespace" = var.k8s_namespace
    }
    "spec" = {
      "entrypoints" = ["web"]
      "routes" = [
        {
          "match" = "Host(`${var.routing.host_match}`)"
          "kind" = "Rule"
          "services" = [
            {
              "name" = var.routing.service_name
              "port" = var.routing.service_port
            }
          ]
        }
      ]
    }
  }
}
