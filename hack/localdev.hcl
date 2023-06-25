# This job is useful for doing local develpment on nomad-holepunch.
#
# To use this job, first 'go install' the project so that the nomad-holepunch
# executable exists on your $PATH.
#
# Then just run the job: nomad job run -var=user=$USER localdev.hcl
#
# The Nomad API will be exposed on port 3030 by default, and is configurable
# with -var=port=<port>

variable "user" {
  type = string
  default = "root"
}

variable "port" {
  type = number
  default = 3030
}

job "localdev" {
  group "group" {
    update {
      min_healthy_time = "2s"
    }

    reschedule {
      attempts  = 0
      unlimited = false
    }

    restart {
      attempts = 0
      mode     = "fail"
    }

    network {
      mode = "host"
      port "api" {
        static = "${var.port}"
      }
    }

    service {
      provider = "nomad"
      port     = "api"
      name     = "holepunch"
      check {
        type     = "http"
        path     = "/health"
        interval = "6s"
        timeout  = "1s"
      }
    }

    task "holepunch" {
      user   = "${var.user}"
      driver = "raw_exec"

      config {
        command = "nomad-holepunch"
      }

      env {
        HOLEPUNCH_PORT = "${NOMAD_PORT_api}"
      }

      identity {
        env = true
      }

      resources {
        cpu    = 100
        memory = 16
      }
    }
  }
}
