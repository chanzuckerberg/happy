---
layout: default
title: Design
nav_order: 3
---
 
# Design
 
There are two main pieces of happy's design:

1. Long-lived infrastructure (aka "happy environments")
1. Short-lived infrastructure (aka "happy stacks")

![both](https://github.com/chanzuckerberg/happy/assets/76011913/333cbfca-8b0e-40f8-84a5-5a49be3d69a1)

## Environment

The long-lived infrastructure is what we generally call the happy environment. It is what you deploy to. It includes the following things:

* EKS or ECS cluster
* Networking
* CI/CD workflows
* Batch processing
* Aurora Postgres DB cluster
* OIDC 
* S3 buckets
* Terraform workspaces
* Any other custom or advanced infra

We refer to it as long-lived instructure because once it is created, it does not change a lot. These are generally resources we don't want
to bring up and down a lot or take a long time to provision. They might have data that can't be destroyed. The long-lived infrastructure
is what we build our stacks on top of.

![happy-env](https://github.com/chanzuckerberg/happy/assets/76011913/46185bae-f3f6-4ffa-a47f-d4497a4bdbac)

## Stacks

The short-lived infrastructure is what we generally call happy stacks. Each stack represents a fully encapsulated version of your application
to deploy. Stacks are:

* A containerized application
* Containerized jobs for batch, cron or other processing
* Runtime configuration (key, value pairs)
* Any small custom infrastructure

Each stack is a complete running application. They are meant to come up and down very quickly. They are great for:

* Testing a new feature on a PR
* Showcasing your work in a demo environment
* Spinning up a pentesting environment
* Spinning up a development environment
* Iteration on microservice locally

![stack](https://github.com/chanzuckerberg/happy/assets/76011913/3a9c4fa3-cc2f-4e0c-b384-d6dee3c0de2c)

## Integration Secret

Happy stacks need to know where to deploy to. In other words, the short-lived infrastructure needs to know about the long-lived infrastructure.
This information is conveyed in a JSON documented called the integration secret. The documents
looks like this:

~~~json
{
    "certificate_arn": "arn:aws:acm:us-west-2:401986845158:certificate/fe87ad86-23b0-41b4-b2ff-e1724a004d59",
    "ci_roles": [
        {
            "arn": "arn:aws:iam::401986845158:role/gh_actions_si_playground_eks_v2",
            "name": "gh_actions_si_playground_eks_v2"
        }
    ],
    "cloud_env": {
        "database_subnet_group": "si-playground",
        "database_subnets": [
            "subnet-00fc2b9327edbc4d9",
            "subnet-0fb4dd8fa2a6e5f18",
            "subnet-0a8964bf1c4d576e4",
            "subnet-0b8efc311727e0255"
        ],
        "private_subnets": [
            "subnet-0ffdf3a2f10adb1f9",
            "subnet-0d69f41508803494c",
            "subnet-08c029e8a4a6b6a32",
            "subnet-0fa1ce01ef263099c"
        ],
        "public_subnets": [
            "subnet-0bfdc8898996f59f2",
            "subnet-0af8d4c628578cd67",
            "subnet-0aa5f4235faa1c475",
            "subnet-0881e4e79011778cc"
        ],
        "vpc_cidr_block": "10.252.0.0/16",
        "vpc_id": "vpc-00c10bee0068a11f2"
    },
    "dbs": {},
    "dynamo_locktable_name": "si-prod-happy-eks-prod-stacklist",
    "ecrs": {},
    "eks_cluster": {
        "cluster_arn": "arn:aws:eks:us-west-2:401986845158:cluster/si-playground-eks-v2",
        "cluster_ca": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMvakNDQWVhZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJek1EVXlOREl3TWpnME1sb1hEVE16TURVeU1USXdNamcwTWxvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTjNYCjduRzRhREQxcHRpb2Q1enJjNlFtZ1YyMnZtVWk5azFWY0tjRDE3WDAzOUk3dUc5S3JXRURxUXlrSGpaYmg0VzQKYzN6YmtFZW9FMlJ0R1p2TGowL1duLzNuazl3ZElETnZUb1pQd3ZqbGtFSVVQVnFrMEc3a3hYREt6M3VnSVBRMgpXQnlHOGJLU0FpUDlneXp3dktYTWVjSkEvYVBpVUhZMkhWRFFPbG5SNkM0dTFqMG9md2I5di9ReHY4Y3NPbE1ECjJIUFh4TVBSODdTZ2UxQzRvelZMTGNzS0Vyc3hYZ0NNd1h6T21xajBHTjg4RVBJSnVFQUQwcHF4aU9nRi90YXUKek92S1VVcTFUMGFCektqVjVqVWZGYXU4aGUyTHBYZzI4cUFKbzlpZGNibXk3WG5PTGhrdjgzdDlRcEVlYjRVegpRcldrZ1ZINzN1Q0g3Q2NoTExVQ0F3RUFBYU5aTUZjd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZLQkdvZWg2di9EejM5WjNVbDFPemJQMkl4Z1BNQlVHQTFVZEVRUU8KTUF5Q0NtdDFZbVZ5Ym1WMFpYTXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSzRCOEZrSVp5cW15N1A3TlZKQQpxNzRoL3lieXpKcUZWU3VrdmVCd0tkSlE1WnEyaU1UeldaQU95K2pqY2xtZkYrMDJBU0xyMWlSRnYwaXV5MEROCjh5V0pja3ZNcWJRTm9raFhoVFltMjBHbmVNaEd6cW9jNUlmTmpKTUlVQzlaSFUzYVRtV3ZQVXFVM0RKd0FQZHQKT0pWMTR3Y3F4WWxMMkJ6dElwOUlpK0ZyNmp5NzFWT25haEphV1VudHNjMnFpTDZ1cE5ZMjk1TmkwV0R2V00vVQoyWmdYTHRNemVvanMxdWdIUWZkb3NuRGo2OUZFUDk4V1VrQzg5NjVxcno4SjV3ZENvVnI2ZnMvWGRPeElDOWtwCnNzVHlldm1Jb1pKalBiQmJNQnhEWHVaMHBSUmxpbnZRd3A0RElPbUdRN0w2WmQvd1pGa2FWOFRxbXdpQXJhdDIKM0lrPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
        "cluster_endpoint": "https://FF38057B0DB0747C32795D4B8265796A.gr7.us-west-2.eks.amazonaws.com",
        "cluster_id": "si-playground-eks-v2",
        "cluster_oidc_issuer_url": "https://oidc.eks.us-west-2.amazonaws.com/id/FF38057B0DB0747C32795D4B8265796A",
        "cluster_version": "1.26",
        "oidc_provider_arn": "arn:aws:iam::401986845158:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/FF38057B0DB0747C32795D4B8265796A",
        "worker_iam_role_name": "si-playground-eks-v2-eks-node20230524200454767700000001",
        "worker_security_group": "sg-0d05e6512e80e6a75"
    },
    "external_zone_name": "happy-playground-prod.prod.si.czi.technology",
    "hapi_config": {
        "assume_role_arn": "arn:aws:iam::401986845158:role/tfe-si",
        "base_url": "https://hapi.hapi.prod.si.czi.technology",
        "kms_key_id": "358d3b31-aa91-4017-b266-214001e37d41",
        "oidc_authz_id": "aus8zrryrcU8fOk9y5d7",
        "oidc_issuer": "0oa8zrudz1DLeTqOG5d7",
        "scope": "happy-eks-prod"
    },
    "kind": "k8s",
    "oidc_config": {
        "client_id": "0oa8zrsvuh6c9hoHe5d7",
        "client_secret": "KAzUaIlYMqHHkbK2AnF1xXhuCsD1j84GL06WwxGg",
        "config_uri": "https://0oa8zrsvuh6c9hoHe5d7:KAzUaIlYMqHHkbK2AnF1xXhuCsD1j84GL06WwxGg@czi-prod.okta.com/oauth2/",
        "idp_url": "czi-prod.okta.com"
    },
    "tags": {
        "env": "prod",
        "managedBy": "terraform",
        "owner": "infra-eng@chanzuckerberg.com",
        "project": "si",
        "service": "happy-eks-prod"
    },
    "tfe": {
        "org": "happy-playground",
        "url": "https://si.prod.tfe.czi.technology"
    },
    "vpc_id": "vpc-00c10bee0068a11f2",
    "zone_id": "Z10432202W0G3BIMH0531"
}
~~~

This is a lot of information in here but the general outline looks like the following:

* VPC information
* Database information (username, host, password)
* EKS/ECS cluster to deploy to 
* Tags
* OIDC information (for protected environments)
* TLS certificate information
* IAM roles 

When a stack is deploying, the stack will have access to the integration secret so that it can deploy in the correct place. It 
can also be used at runtime too. When a stack is deployed, the integration secret will be mounted to the deployed container
at `/var/happy/integration-secret`. Applications can use this to connect to databases, assume roles or parse other advanced
infrastructure that might be needed at runtime by the application.

The integration secret is a very handy contract that is agreed upon between the short and long lived infra. This allows
happy developers to change their infra or bring their own infra if they don't want to use the happy defaults or if they
have more advanced use cases. As long as they implement the integration secret, happy stacks can be built on top of it.
Developers can also add fields to the integration secret as needed to include these more advanced infrastructure requirements.

***NOTE***: happy looks for the integration secret in a specific location -- a kubernetes secret called "integration-secret"
under the namespace of your happy application. Here's an example:

~~~yaml
apiVersion: v1
kind: Secret
metadata:
  name: integration-secret
  namespace: edu-platform-rdev-happy-happy-env
data:
  integration_secret: >-
    eyJjZXJ0aWZpY2F0ZV9hcm4iOiJhcm46YXdzOmFjbTp1cy13ZXN0LTI6MjU0Mzc5NDUzNzA3OmNlcnRpZmljYXRlLzljYzNkODQwLTkwMDEtNDkwZi05M2JmLTY3NTUxNjgxNjAxNSIsImNpX3JvbGVzIjpbeyJhcm4iOiJhcm46YXdzOmlhbTo6MjU0Mzc5NDUzNzA3OnJvbGUvZ2hfYWN0aW9uc19lZHVfcGxhdGZvcm1fcmRldl9la3MiLCJuYW1lIjoiZ2hfYWN0aW9uc19lZHVfcGxhdGZvcm1fcmRldl9la3MifV0sImNsb3VkX2VudiI6eyJkYXRhYmFzZV9zdWJuZXRfZ3JvdXAiOiJlZHUtcGxhdGZvcm0tcmRldiIsImRhdGFiYXNlX3N1Ym5ldHMiOlsic3VibmV0LTA1YzFkNDEyYzBlMGQxZDVhIiwic3VibmV0LTAyNGE5NzkxOGJiNmJjM2NjIiwic3VibmV0LTBjMzQ3NDI2MmE0MTZkNTAzIiwic3VibmV0LTBmNDQwNjRhYjk1Njk0NDYwIl0sInByaXZhdGVfc3VibmV0cyI6WyJzdWJuZXQtMDk0Yjc4MGQzZWI2ZGUyYTkiLCJzdWJuZXQtMDkwMDY3Zjc1MWJhYjlkMDYiLCJzdWJuZXQtMDllMTM1YzJkNGJmMTg0MTMiLCJzdWJuZXQtMDUzM2Q1ZTY3ZTllMDA5NzIiXSwicHVibGljX3N1Ym5ldHMiOlsic3VibmV0LTAyNTVjZjlmMDlkMTY2MWIxIiwic3VibmV0LTBhYzZjYTUxZWE0OGE4MzJkIiwic3VibmV0LTA2MGJhMWI3NTljNWViZWFjIiwic3VibmV0LTBlZDc0ZDAxMmQ4OWE0MTNkIl0sInZwY19jaWRyX2Jsb2NrIjoiMTAuOTkuMC4wLzE2IiwidnBjX2lkIjoidnBjLTAxN2I0ODc5NTNlZTdlZDVlIn0sImRicyI6eyJzZGwiOnsiZGF0YWJhc2VfaG9zdCI6ImVkdS1wbGF0Zm9ybS1yZGV2LWhhcHB5LXNkbC5jbHVzdGVyLWNpdTJ3aXlxazE4Zi51cy13ZXN0LTIucmRzLmFtYXpvbmF3cy5jb20iLCJkYXRhYmFzZV9uYW1lIjoic2RsIiwiZGF0YWJhc2VfcGFzc3dvcmQiOiJmcG51VUhkcWlpaUI2Y1pGS1VnaGIyWkVnekVhMHQ1NyIsImRhdGFiYXNlX3BvcnQiOiI1NDMyIiwiZGF0YWJhc2VfdXNlciI6InNkbCJ9fSwiZHluYW1vX2xvY2t0YWJsZV9uYW1lIjoiZWR1LXBsYXRmb3JtLXJkZXYtaGFwcHktc3RhY2tsaXN0IiwiZWNycyI6eyJjZGMtZGViZXppdW0tY29uc3VtZXIiOnsidXJsIjoiMjU0Mzc5NDUzNzA3LmRrci5lY3IudXMtd2VzdC0yLmFtYXpvbmF3cy5jb20vY2RjLWRlYmV6aXVtLWNvbnN1bWVyIn19LCJla3NfY2x1c3RlciI6eyJjbHVzdGVyX2FybiI6ImFybjphd3M6ZWtzOnVzLXdlc3QtMjoyNTQzNzk0NTM3MDc6Y2x1c3Rlci9lZHUtcGxhdGZvcm0tcmRldi1la3MiLCJjbHVzdGVyX2NhIjoiTFMwdExTMUNSVWRKVGlCRFJWSlVTVVpKUTBGVVJTMHRMUzB0Q2sxSlNVTXZha05EUVdWaFowRjNTVUpCWjBsQ1FVUkJUa0puYTNGb2EybEhPWGN3UWtGUmMwWkJSRUZXVFZKTmQwVlJXVVJXVVZGRVJYZHdjbVJYU213S1kyMDFiR1JIVm5wTlFqUllSRlJKZWsxRVNYbE9SRUYzVFhwTmVVOVdiMWhFVkUxNlRVUkplVTFVUVhkTmVrMTVUMVp2ZDBaVVJWUk5Ra1ZIUVRGVlJRcEJlRTFMWVROV2FWcFlTblZhV0ZKc1kzcERRMEZUU1hkRVVWbEtTMjlhU1doMlkwNUJVVVZDUWxGQlJHZG5SVkJCUkVORFFWRnZRMmRuUlVKQlRHOTJDbkZCVHlzNE5GRm5VMlYzVUdabVpXY3liRmR5YnpodU5EQjJXbXhyTUV0aVRVNUVkV015UmtkQmRWUTFiRXBCVW5WQ01FRkpZblF4VDJ4YVMyeEVNaXNLWkVGMVlsQnZhSFpHV2tSck56TlplbXczTkhOa2JETlpUbFZCU2xodWNtTnNWbEZqVDJwelZXOTFSbk5VU0ZjMFptRlhTazk0U0dWeGNWUXlOWFIzTVFwVVNsRklhSFFyZGpBNGEwTlBTSGNyZVZkM2NqSkZhWGt6TkhOT1RrZ3JXRnBVWVhGbUswRlNTbVJ4YWtrd1pVbHlSU3RDTUZwMGJVVXdZazlSVG1WdkNraEZRa3c0U1dOTVJYTmhUMnAxYkVoaWFFbGtaamxIVERSWVVrVm9lRTFsVVdkT2FHaEJSRXQ2T1ZGTVozWXJiVlJZVEdNd0szUnhlR041Y0VKdFpuSUtNWE16T0ZSMFRGbGxOelkzYTJGMFYyTjZNemgxWjFVM1ZXbG5PVTl4WWxsU1JuQmhaeTlIWTBkWVRqRlRTSFU0VGxVMWRuRndiRWhVYlRBeUt5OXZiZ292WjJ4T2NVYzBTVGgzYjJGcFZqZGhNMlpWUTBGM1JVRkJZVTVhVFVaamQwUm5XVVJXVWpCUVFWRklMMEpCVVVSQlowdHJUVUU0UjBFeFZXUkZkMFZDQ2k5M1VVWk5RVTFDUVdZNGQwaFJXVVJXVWpCUFFrSlpSVVpMY0RrMU1IZDViR0l4YjJ4NFdFOHhkVWQ2VGpSeFNIZzFORnBOUWxWSFFURlZaRVZSVVU4S1RVRjVRME50ZERGWmJWWjVZbTFXTUZwWVRYZEVVVmxLUzI5YVNXaDJZMDVCVVVWTVFsRkJSR2RuUlVKQlJYbHVjVGRCZW05cmRIRnhVamxMVWpSdUt3cDNUVVY2WlRKUmVsSkpSVEJRVGxjNFExbHJUSEZrZGtRMmJreFpVMUUxVWsxV1JXRTFTa0ppUjJwcWNYaGtUbHAwVWtWSVkyaG9ZbVl3VTJoRWNpdHRDbEJ5WW01bVdEUjBNVWgwZDJGbFkzTTNhMWh6U1ZBcmRUVm9iRkp3VjFKb1YxQktjMDVuWTIxaksxWkxTRkF3ZDBGalUzcG1SR2cyWmtzeGMzQkhTbXdLVDIxWE5UTmhTelkzU0M5VWVub3JPV0pNUlZWVlJuUlJNbEpvVFZWNlkwOUZXa0ZZVkV4R1ZtVnRja3RtV0c5TmNDdFNOM1pVSzJsMlpESklaWGRsTlFvNFptcFZabVpHU1UxbWNuZFdTV1UwVlZaemJFOTVRV0ZMVjNaRWFtOXdlbXBCY1ZCWWN6RnZVRWs1YUZFcmFXOVNTMVJuUTBnNFVGWTFSalV5TVd0M0NtTTJZVTV6V1ZjMU5tMWlVMGs1TmpCMVpWRm5TMjgzVVd0Q1pFcGtjazUzUzFwTVpWRjVNbFpsZUhCRE0zTjRhMDFFYjNGcmIwaGFVME5YYm5aMGFXY0tkV1ZaUFFvdExTMHRMVVZPUkNCRFJWSlVTVVpKUTBGVVJTMHRMUzB0Q2c9PSIsImNsdXN0ZXJfZW5kcG9pbnQiOiJodHRwczovLzQzNUE1MUIxRjYzNDgwMDlGN0ZCODBGNDIyNzY5QjFFLmdyNy51cy13ZXN0LTIuZWtzLmFtYXpvbmF3cy5jb20iLCJjbHVzdGVyX2lkIjoiZWR1LXBsYXRmb3JtLXJkZXYtZWtzIiwiY2x1c3Rlcl9vaWRjX2lzc3Vlcl91cmwiOiJodHRwczovL29pZGMuZWtzLnVzLXdlc3QtMi5hbWF6b25hd3MuY29tL2lkLzQzNUE1MUIxRjYzNDgwMDlGN0ZCODBGNDIyNzY5QjFFIiwiY2x1c3Rlcl92ZXJzaW9uIjoiMS4yNSIsIm9pZGNfcHJvdmlkZXJfYXJuIjoiYXJuOmF3czppYW06OjI1NDM3OTQ1MzcwNzpvaWRjLXByb3ZpZGVyL29pZGMuZWtzLnVzLXdlc3QtMi5hbWF6b25hd3MuY29tL2lkLzQzNUE1MUIxRjYzNDgwMDlGN0ZCODBGNDIyNzY5QjFFIiwid29ya2VyX2lhbV9yb2xlX25hbWUiOiItZWtzLW5vZGUyMDIzMDIyNDAwMjcyMDEyNjQwMDAwMDAwMSIsIndvcmtlcl9zZWN1cml0eV9ncm91cCI6InNnLTAxODRmNTI1MTFlMTRjODJhIn0sImV4dGVybmFsX3pvbmVfbmFtZSI6ImVkdS1wbGF0Zm9ybS5yZGV2LnNpLmN6aS50ZWNobm9sb2d5IiwiaGFwaV9jb25maWciOnsiYXNzdW1lX3JvbGVfYXJuIjoiYXJuOmF3czppYW06OjI1NDM3OTQ1MzcwNzpyb2xlL3RmZS1zaSIsImJhc2VfdXJsIjoiaHR0cHM6Ly9oYXBpLmhhcGkucHJvZC5zaS5jemkudGVjaG5vbG9neSIsImttc19rZXlfaWQiOiI4NGQ1NjYwOC02ZGU3LTQyNzctYjFkNi0xZDI5ZWEyOWRiYjMiLCJvaWRjX2F1dGh6X2lkIjoiYXVzOGY0OWEwdlliN2RIZmU1ZDciLCJvaWRjX2lzc3VlciI6IjBvYThmNDRva2MzdjBHQkFvNWQ3Iiwic2NvcGUiOiJoYXBweSJ9LCJraW5kIjoiazhzIiwib2lkY19jb25maWciOnsiY2xpZW50X2lkIjoiMG9hODY5MGQ0dTdMOVBEdEU1ZDciLCJjbGllbnRfc2VjcmV0IjoiTjQxQUVRSGE4NTFtRFRJNWJ4QnlrMnZHZEJzRWo0VldjNDZsVGt5TCIsImNvbmZpZ191cmkiOiJodHRwczovLzBvYTg2OTBkNHU3TDlQRHRFNWQ3Ok40MUFFUUhhODUxbURUSTVieEJ5azJ2R2RCc0VqNFZXYzQ2bFRreUxAY3ppLXByb2Qub2t0YS5jb20vb2F1dGgyLyIsImlkcF91cmwiOiJjemktcHJvZC5va3RhLmNvbSJ9LCJ0YWdzIjp7ImVudiI6InJkZXYiLCJtYW5hZ2VkQnkiOiJ0ZXJyYWZvcm0iLCJvd25lciI6ImluZnJhLWVuZ0BjaGFuenVja2VyYmVyZy5jb20iLCJwcm9qZWN0IjoiZWR1LXBsYXRmb3JtIiwic2VydmljZSI6ImhhcHB5In0sInRmZSI6eyJvcmciOiJoYXBweS1lZHUtcGxhdGZvcm0iLCJ1cmwiOiJodHRwczovL3NpLnByb2QudGZlLmN6aS50ZWNobm9sb2d5In0sInZwY19pZCI6InZwYy0wMTdiNDg3OTUzZWU3ZWQ1ZSIsIndhZl9jb25maWciOnsiYXJuIjoiYXJuOmF3czp3YWZ2Mjp1cy13ZXN0LTI6MjU0Mzc5NDUzNzA3OnJlZ2lvbmFsL3dlYmFjbC9lZHUtcGxhdGZvcm0tcmRldi1oYXBweS82OWNjYmRlOS02NjI2LTQ0YzQtYjcxZC00MzBhYTdiMGVlM2UiLCJzY29wZSI6IlJFR0lPTkFMIn0sInpvbmVfaWQiOiJaMTAxMzY1MFJFMVg5TUZPR1pUOCJ9
type: Opaque
~~~

# References:

* [Presentation](https://docs.google.com/presentation/d/1zgbTF_1oq96npmKXxHKFVn5rO96wEsQlj5bgd7axLNA/edit#slide=id.p)
* [Stack options](https://github.com/chanzuckerberg/happy/blob/0142e747802df4768f1d27e27a062f86a821316d/terraform/modules/happy-stack-eks/variables.tf#L46)