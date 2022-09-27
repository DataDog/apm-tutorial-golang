data "external" "whatismyip" {
  program = ["/bin/bash" , "../whatismyip.sh"]
}

# Providing a reference to our default VPC
resource "aws_default_vpc" "default_vpc" {
  enable_dns_support = true
  enable_dns_hostnames = true
}

# Providing a reference to our default subnets
resource "aws_default_subnet" "default_subnet_a" {
  availability_zone = "us-east-1a"
}

resource "aws_default_subnet" "default_subnet_b" {
  availability_zone = "us-east-1b"
}

resource "aws_default_subnet" "default_subnet_c" {
  availability_zone = "us-east-1c"
}