# terraform-provider-cognito-idp-link

This Terraform Provider provides functionality to link federated users to existing user profiles in Amazon Cognito User Pools.

## Features

- Link federated users to existing user profiles
- Compatible with various identity providers including SAML and OIDC

## Usage

This provider can be used alongside the official AWS provider. Here's an example of how to use both providers together:

```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    cognito-idp-link = {
      source = "atpons/cognito-idp-link"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"
}

provider "cognito-idp-link" {}

# Create a Cognito User Pool using AWS provider
resource "aws_cognito_user_pool" "example" {
  name = "example-pool"
}

# Link federated users using this provider
resource "cognito-idp-link_link" "example" {
  user_pool_id = aws_cognito_user_pool.example.id

  destination_user = {
    provider_name            = "Cognito"
    provider_attribute_name  = "userId"
    provider_attribute_value = "87f45a38-1091-70ca-a0f9-4e2110859e84"
  }

  source_user = {
    provider_name            = "oidc"
    provider_attribute_name  = "email"
    provider_attribute_value = "test@example.com"
  }
}
```
