resource "cognito_idp_link" "example" {
  user_pool_id = "ap-northeast-1_XXXXXXXX"

  destination_user = {
    provider_name            = "Cognito"
    provider_attribute_name  = "userId"
    provider_attribute_value = "87f45a38-1091-70ca-a0f9-4e2110859e84"
  }

  source_user {
    provider_name            = "oidc"
    provider_attribute_name  = "email"
    provider_attribute_value = "admin@example.com"
  }
}
