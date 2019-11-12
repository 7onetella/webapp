# This declares a job named "docs". There can be exactly one
# job declaration per job file.

job "webapp" {
  # Spread the tasks in this job between us-west-1 and us-east-1.
  datacenters = ["dc1"]

  # Run this job as a "service" type. Each job type has different
  # properties. See the documentation below for more examples.
  type = "service"

  # Specify this job to have rolling updates, two-at-a-time, with
  # 30 second intervals.
  update {
    stagger      = "30s"
    max_parallel = 1
  }

  # A group defines a series of tasks that should be co-located
  # on the same client (host). All tasks within a group will be
  # placed on the same host.
  group "web" {
    # Specify the number of these tasks we want.
    count = 2

    # Create an individual task (unit of work). This particular
    # task utilizes a Docker container to front a web application.
    task "app" {
      artifact {
        source      = "http://vagrant.7onetella.net/uploads/webapp_BUILD_ID.tar.gz"
      }

      # Specify the driver to be "docker". Nomad supports
      # multiple drivers.
      driver = "raw_exec"

      # Configuration is specific to each driver.
      config {
        command = "webapp_BUILD_ID"
      }

      # The service block tells Nomad how to register this service
      # with Consul for service discovery and monitoring.
      service {
        tags = ["urlprefix-webapp.7onetella.net/"]
        
        # This tells Consul to monitor the service on the port
        # labelled "http". Since Nomad allocates high dynamic port
        # numbers, we use labels to refer to them.
        port = "http"

        check {
         type     = "http"
         path     = "/api/health"
         interval = "10s"
         timeout  = "2s"
        }
      }

      # It is possible to set environment variables which will be
      # available to the task when it runs.
      env {
        "HTTP_PORT" = "${NOMAD_PORT_http}"
      }

      template {
        data = <<EOH
BUILD_NUMBER="${BUILD_NUMBER}"
EOH
        destination = "secrets/file.env"
        env         = true
      }

      # Specify the maximum resources required to run the task,
      # include CPU, memory, and bandwidth.
      resources {
        cpu    = 128 # MHz
        memory = 128 # MB

        network {
          mbits = 100
          port "http" {}
        }
      }
    }
  }
}
