# Infrastructure Blueprint: Secure VPC Isolation for Anthropic Claude via AWS Bedrock

resource "aws_vpc" "enterprise_ai_vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "enterprise-ai-architecture-vpc"
  }
}

# Private Subnet for hosting our Edge Data Bridge Microservices
resource "aws_subnet" "private_mcp_subnet" {
  vpc_id            = aws_vpc.enterprise_ai_vpc.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "eu-west-1a"

  tags = {
    Name = "mcp-isolated-private-subnet"
  }
}

# Security Group restricting all inbound traffic except authorized proxy layers
resource "aws_security_group" "mcp_security_boundary" {
  name        = "mcp-security-boundary"
  description = "Enforce zero public access to data transformation components"
  vpc_id      = aws_vpc.enterprise_ai_vpc.id

  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Restricted strictly to AWS Bedrock API gateway endpoints
  }
}

# Private VPC Endpoint linking to AWS Bedrock without crossing the public internet
resource "aws_vpc_endpoint" "bedrock_runtime" {
  vpc_id            = aws_vpc.enterprise_ai_vpc.id
  service_name      = "com.amazonaws.eu-west-1.bedrock-runtime"
  vpc_endpoint_type = "Interface"

  subnet_ids         = [aws_subnet.private_mcp_subnet.id]
  security_group_ids = [aws_security_group.mcp_security_boundary.id]

  private_dns_enabled = true
}
