# Creating a security group for the load balancer:
resource "aws_security_group" "load_balancer_security_group" {
  vpc_id             = "${aws_default_vpc.default_vpc.id}"
  ingress {
    from_port   = 0
    to_port     = 10000
    protocol    = "tcp"
    cidr_blocks = [format("%s/%s", data.external.whatismyip.result["internet_ip"],32)]
    //Replace the below value with machine's own public IPv4 address if necessary
    //cidr_blocks = ["127.0.0.1/32"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "service_security_group" {
  vpc_id             = "${aws_default_vpc.default_vpc.id}"
  ingress {
    from_port = 0
    to_port   = 9090
    protocol  = "tcp"
    # Only allowing traffic in from the VPC 
    cidr_blocks = [aws_default_vpc.default_vpc.cidr_block]
  }

  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    # Only allowing traffic in from the load balancer security group
    security_groups = ["${aws_security_group.load_balancer_security_group.id}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}