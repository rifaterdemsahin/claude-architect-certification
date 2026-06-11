# ==============================================================================
# AWS Terraform Blueprint: VPC Isolation & Bedrock Private Link for Claude
# ==============================================================================

terraform {
  required_version = ">= 1.3.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

# 1. Isolated Virtual Private Cloud (VPC)
resource "aws_vpc" "secure_ml_vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name        = "secure-ml-vpc"
    Environment = "production"
    Compliance  = "zero-data-retention"
  }
}

# 2. Private Subnet (No Internet Gateway routing)
resource "aws_subnet" "private_subnet" {
  vpc_id            = aws_vpc.secure_ml_vpc.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "${var.aws_region}a"

  tags = {
    Name = "secure-ml-private-subnet"
  }
}

# 3. Security Group controlling egress traffic
resource "aws_security_group" "mcp_security_group" {
  name        = "mcp-server-sg"
  description = "Restrict outgoing traffic to Bedrock interface endpoints only"
  vpc_id      = aws_vpc.secure_ml_vpc.id

  # Inbound connection rules (e.g., internal corp proxy/gateway)
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }

  # Egress rules: Only allow HTTPS to VPC endpoints inside the CIDR range
  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }
}

# 4. VPC Endpoint for Amazon Bedrock Runtime (accessing Claude)
resource "aws_vpc_endpoint" "bedrock_runtime_endpoint" {
  vpc_id            = aws_vpc.secure_ml_vpc.id
  service_name      = "com.amazonaws.${var.aws_region}.bedrock-runtime"
  vpc_endpoint_type = "Interface"

  subnet_ids         = [aws_subnet.private_subnet.id]
  security_group_ids = [aws_security_group.mcp_security_group.id]
  private_dns_enabled = true

  tags = {
    Name = "bedrock-runtime-vpce"
  }
}
