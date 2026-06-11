# 📐 templates — Infrastructure Templates

> **Purpose:** Terraform and infrastructure-as-code templates for cloud resource provisioning.

## Files

| File | Description |
|------|-------------|
| `aws_bedrock_vpc_endpoint.tf` | Terraform config for AWS Bedrock VPC endpoint |
| `aws-bedrock-private-link.tf` | Terraform config for AWS Bedrock PrivateLink |

## Rules
- Never commit live credentials in `.tfvars` files
- Reference secrets from Azure Key Vault, not plaintext