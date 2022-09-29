module "settings" {
  source = "../global_constants"
}

provider "aws" {
  region  = module.settings.aws_region
  profile = module.settings.aws_profile
}

resource "aws_ecs_cluster" "apm_tutorial_cluster" {
  name = "apm-tutorial-fargate-go" # Naming the cluster
}

resource "aws_ecs_task_definition" "notes_task" {
  family = "notes-task" # Naming our first task
  container_definitions = jsonencode([
    {
      name : "notes-task",
      image : "${module.settings.aws_ecr_repository}:notes",
      essential : true,
      portMappings : [
        {
          containerPort : 8080,
          hostPort : 8080
        }
      ],
      memory : 512,
      cpu : 256,
      environment : [
        {
          name : "CALENDAR_HOST",
          value : "calendar.apmlocalgo"
        },
        {
          name : "DD_SERVICE",
          value : "notes"
        },
        {
          name : "DD_ENV",
          value : "dev"
        },
        {
          name : "DD_VERSION",
          value : "0.0.1"
        }
      ],
      dockerLabels : {
        "com.datadoghq.tags.service" : "notes",
        "com.datadoghq.tags.env" : "dev",
        "com.datadoghq.tags.version" : "0.0.1"
      },
    },
    {
      name : "datadog-agent",
      image : "public.ecr.aws/datadog/agent:latest",
      essential : true,
      environment : [
        {
          name : "DD_API_KEY",
          value : module.settings.datadog_api_key
        },
        {
          name : "ECS_FARGATE",
          value : "true"
        },
        {
          name : "DD_APM_ENABLED",
          value : "true"
        }
      ],
      portMappings : [
        {
          containerPort : 8126,
          protocol : "tcp",
        }
      ],
    }
  ])
  requires_compatibilities = ["FARGATE"] # Stating that we are using ECS 
  network_mode             = "awsvpc"    # Using awsvpc as our network mode as this is required
  memory                   = 512         # Specifying the memory our container requires
  cpu                      = 256         # Specifying the CPU our container requires
  execution_role_arn       = aws_iam_role.ecsTaskExecutionRole.arn
}

resource "aws_ecs_service" "notes_service" {
  name            = "notes-service"                         # Naming our first service
  cluster         = aws_ecs_cluster.apm_tutorial_cluster.id # Referencing our created Cluster
  task_definition = aws_ecs_task_definition.notes_task.arn  # Referencing the task our service will spin up
  launch_type     = "FARGATE"
  desired_count   = 1 # Setting the number of containers to 3

  load_balancer {
    target_group_arn = aws_lb_target_group.target_group.arn # Referencing our target group
    container_name   = aws_ecs_task_definition.notes_task.family
    container_port   = 8080 # Specifying the container port
  }

  service_registries {
    registry_arn   = aws_service_discovery_service.apm_notes_service.arn
    container_name = "notes"
  }

  network_configuration {
    subnets          = ["${aws_default_subnet.default_subnet_a.id}", "${aws_default_subnet.default_subnet_b.id}", "${aws_default_subnet.default_subnet_c.id}"]
    assign_public_ip = true                                                # Providing our containers with public IPs
    security_groups  = ["${aws_security_group.service_security_group.id}"] # Setting the security group
  }
}

## ----------------------------------------------- Calendar Service ------------------------------------------------------------------------------------------


resource "aws_ecs_task_definition" "calendar_task" {
  family = "calendar-task" # Naming our first task
  container_definitions = jsonencode([
    {
      name : "calendar-task",
      image : "${module.settings.aws_ecr_repository}:calendar",
      essential : true,
      environment : [
        {
          name : "DD_SERVICE",
          value : "calendar"
        },
        {
          name : "DD_ENV",
          value : "dev"
        },
        {
          name : "DD_VERSION",
          value : "0.0.1"
        }
      ],
      dockerLabels : {
        "com.datadoghq.tags.service" : "calendar",
        "com.datadoghq.tags.env" : "dev",
        "com.datadoghq.tags.version" : "0.0.1"
      },
      portMappings : [
        {
          containerPort : 9090,
          hostPort : 9090
        }
      ],
      memory : 512,
      cpu : 256,
    },
    {
      name : "datadog-agent",
      image : "public.ecr.aws/datadog/agent:latest",
      essential : true,
      environment : [
        {
          name : "DD_API_KEY",
          value : module.settings.datadog_api_key
        },
        {
          name : "ECS_FARGATE",
          value : "true"
        },
        {
          name : "DD_APM_ENABLED",
          value : "true"
        }
      ],
      portMappings : [
        {
          containerPort : 8126,
          protocol : "tcp",
        }
      ],
    }
  ])
  requires_compatibilities = ["FARGATE"] # Stating that we are using ECS 
  network_mode             = "awsvpc"    # Using awsvpc as our network mode as this is required
  memory                   = 512         # Specifying the memory our container requires
  cpu                      = 256         # Specifying the CPU our container requires
  execution_role_arn       = aws_iam_role.ecsTaskExecutionRole.arn
}

resource "aws_ecs_service" "calendar_service" {
  name            = "calendar-service"                        # Naming our first service
  cluster         = aws_ecs_cluster.apm_tutorial_cluster.id   # Referencing our created Cluster
  task_definition = aws_ecs_task_definition.calendar_task.arn # Referencing the task our service will spin up
  launch_type     = "FARGATE"
  desired_count   = 1 # Setting the number of containers to 1

  load_balancer {
    target_group_arn = aws_lb_target_group.target_group_2.arn # Referencing our target group
    container_name   = aws_ecs_task_definition.calendar_task.family
    container_port   = 9090 # Specifying the container port
  }

  service_registries {
    registry_arn   = aws_service_discovery_service.apm_calendar_service.arn
    container_name = "calendar"
  }

  network_configuration {
    subnets          = ["${aws_default_subnet.default_subnet_a.id}", "${aws_default_subnet.default_subnet_b.id}", "${aws_default_subnet.default_subnet_c.id}"]
    assign_public_ip = true                                                # Providing our containers with public IPs
    security_groups  = ["${aws_security_group.service_security_group.id}"] # Setting the security group
  }
}
