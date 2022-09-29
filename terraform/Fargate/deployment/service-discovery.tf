# dns
resource "aws_service_discovery_private_dns_namespace" "apm_tutorial_dns" {
  name        = "apmlocalgo"
  description = "DNS Space for apm tutorial"
  vpc         = aws_default_vpc.default_vpc.id
}

# calendar service discovery
resource "aws_service_discovery_service" "apm_calendar_service" {
  name = "calendar"

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.apm_tutorial_dns.id

    dns_records {
      ttl  = 10
      type = "A"
    }

    routing_policy = "MULTIVALUE"
  }

  health_check_custom_config {
    failure_threshold = 1
  }
}

# notes service discovery
resource "aws_service_discovery_service" "apm_notes_service" {
  name = "notes"

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.apm_tutorial_dns.id

    dns_records {
      ttl  = 10
      type = "A"
    }

    routing_policy = "MULTIVALUE"
  }

  health_check_custom_config {
    failure_threshold = 1
  }
}
