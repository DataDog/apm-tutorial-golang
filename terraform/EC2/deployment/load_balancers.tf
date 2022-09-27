resource "aws_iam_role" "ecsTaskExecutionRole" {
  name               = "ecsTaskExecutionRoleTfEc2"
  assume_role_policy = "${data.aws_iam_policy_document.assume_role_policy.json}"
}

data "aws_iam_policy_document" "assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy_attachment" "ecsTaskExecutionRole_policy" {
  role       = "${aws_iam_role.ecsTaskExecutionRole.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# -------------------------------------------------- Notes Service LB ------------------------------------------ |

resource "aws_alb" "application_load_balancer" {
  name               = "test-lb-tf-ec2-go" # Naming our load balancer
  load_balancer_type = "application"
  subnets = [ # Referencing the default subnets
    "${aws_default_subnet.default_subnet_a.id}",
    "${aws_default_subnet.default_subnet_b.id}",
    "${aws_default_subnet.default_subnet_c.id}"
  ]
  # Referencing the security group
  security_groups = ["${aws_security_group.load_balancer_security_group.id}"]
}

resource "aws_lb_target_group" "target_group" {
  name        = "apm-tutorial-target-group-ec2"
  port        = 80
  protocol    = "HTTP"
  target_type = "instance"
  vpc_id      = "${aws_default_vpc.default_vpc.id}" # Referencing the default VPC

  health_check {
    matcher = "200"
    path    = "/notes"
    interval = 120
    timeout = 30
  }

}

resource "aws_lb_listener" "listener" {
  load_balancer_arn = "${aws_alb.application_load_balancer.arn}" # Referencing our load balancer
  port              = "8080"
  protocol          = "HTTP"
  default_action {
    type             = "forward"
    target_group_arn = "${aws_lb_target_group.target_group.arn}" # Referencing our target group
  }
}

# -------------------------------------------------- Calendar Service LB ------------------------------------------ |

resource "aws_alb" "application_load_balancer_2" {
  name               = "test-lb-tf-2-ec2-go" # Naming our load balancer
  load_balancer_type = "application"
  subnets = [ # Referencing the default subnets
    "${aws_default_subnet.default_subnet_a.id}",
    "${aws_default_subnet.default_subnet_b.id}",
    "${aws_default_subnet.default_subnet_c.id}"
  ]
  # Referencing the security group
  security_groups = ["${aws_security_group.load_balancer_security_group.id}"]
}

resource "aws_lb_target_group" "target_group_2" {
  name        = "apm-tutorial-target-group-2-ec2"
  port        = 80
  protocol    = "HTTP"
  target_type = "instance"
  vpc_id      = "${aws_default_vpc.default_vpc.id}" # Referencing the default VPC

  health_check {
    matcher = "200"
    path    = "/calendar"
    interval = 120
    timeout = 30
  }

}

resource "aws_lb_listener" "listener_2" {
  load_balancer_arn = "${aws_alb.application_load_balancer_2.arn}" # Referencing our load balancer
  port              = "9090"
  protocol          = "HTTP"
  default_action {
    type             = "forward"
    target_group_arn = "${aws_lb_target_group.target_group_2.arn}" # Referencing our target group
  }
}