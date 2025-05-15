# terraform {
#   required_providers {
#     random = {
#       source  = "hashicorp/random"
#       version = "3.6.3"
#     }
#     local = {
#       source  = "hashicorp/local"
#       version = "2.5.2"
#     }
#     archive = {
#       source  = "hashicorp/archive"
#       version = "2.7.0"
#     }
#   }
# }
#
# provider "random" {
#   # Configuration options
# }
#
# provider "local" {
#   # Configuration options
# }
#
# provider "archive" {
#   # Configuration options
# }
terraform {
  backend "s3" {
    endpoints = {
      s3 = "https://5a8f59dac8e51dbb2086b06288148e94.r2.cloudflarestorage.com"
    }
    bucket = "terraform-state"
    key    = "sausage-factory/baftfexec.tfstate"
    use_lockfile = true

    skip_credentials_validation = true
    skip_region_validation      = true
    skip_metadata_api_check     = true
    skip_requesting_account_id  = true
    skip_s3_checksum = true

    region = "auto"


  }
  required_providers {
    github = {
      source  = "integrations/github"
      version = "6.6.0"
    }
  }

  required_version = "~>1.0"
}

provider "github" {
  # Configuration options
}
