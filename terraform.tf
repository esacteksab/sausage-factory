terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "3.6.3"
    }
    local = {
      source  = "hashicorp/local"
      version = "2.5.2"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "2.7.0"
    }
  }
}

provider "random" {
  # Configuration options
}

provider "local" {
  # Configuration options
}

provider "archive" {
  # Configuration options
}
