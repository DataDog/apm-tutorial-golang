module "settings" {
  source = "../global_constants"
}

provider "aws" {
  region  = module.settings.aws_region
  profile = module.settings.aws_profile
}

resource "aws_ecs_cluster" "apm_tutorial_cluster" {
  name = "apm-tutorial-ec2-go" # Naming the cluster
}

resource "aws_ecs_task_definition" "notes_task" {
  family = "notes" # Naming our first task
  container_definitions = jsonencode([
    {
      name : "notes",
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
          value : "localhost"
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
    }
  ])
  requires_compatibilities = ["EC2"] # Stating that we are using ECS 
  network_mode             = "host"
  memory                   = 512 # Specifying the memory our container requires
  cpu                      = 256 # Specifying the CPU our container requires
  execution_role_arn       = aws_iam_role.ecsTaskExecutionRole.arn
}

resource "aws_ecs_service" "notes_service" {
  name            = "notes-service"                         # Naming our first service
  cluster         = aws_ecs_cluster.apm_tutorial_cluster.id # Referencing our created Cluster
  task_definition = aws_ecs_task_definition.notes_task.arn  # Referencing the task our service will spin up
  launch_type     = "EC2"
  desired_count   = 1 # Setting the number of containers to 3

  load_balancer {
    target_group_arn = aws_lb_target_group.target_group.arn # Referencing our target group
    container_name   = aws_ecs_task_definition.notes_task.family
    container_port   = 8080 # Specifying the container port
  }
}

## ----------------------------------------------- Calendar Service ------------------------------------------------------------------------------------------


resource "aws_ecs_task_definition" "calendar_task" {
  family = "calendar" # Naming our first task
  container_definitions = jsonencode([
    {
      name : "calendar",
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
    }
  ])
  requires_compatibilities = ["EC2"] # Stating that we are using ECS 
  network_mode             = "host"
  memory                   = 512 # Specifying the memory our container requires
  cpu                      = 256 # Specifying the CPU our container requires
  execution_role_arn       = aws_iam_role.ecsTaskExecutionRole.arn
}

resource "aws_ecs_service" "calendar_service" {
  name            = "calendar-service"                        # Naming our first service
  cluster         = aws_ecs_cluster.apm_tutorial_cluster.id   # Referencing our created Cluster
  task_definition = aws_ecs_task_definition.calendar_task.arn # Referencing the task our service will spin up
  launch_type     = "EC2"
  desired_count   = 1 # Setting the number of containers to 1

  load_balancer {
    target_group_arn = aws_lb_target_group.target_group_2.arn # Referencing our target group
    container_name   = aws_ecs_task_definition.calendar_task.family
    container_port   = 9090 # Specifying the container port
  }
}
