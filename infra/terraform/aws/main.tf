terraform {
  backend "s3" {
    # Bucket name, key, and region will be passed via -backend-config during terraform init
  }
}
provider "aws" {
  region = var.region
}

resource "aws_instance" "cd_engine_vm" {
  ami           = var.ami_id
  instance_type = var.instance_type
  key_name      = var.key_name

  root_block_device {
    volume_size = 30
    volume_type = "gp3"
  }

  vpc_security_group_ids = [aws_security_group.cd_engine_sg.id]

  tags = {
    Name = "cd-engine-vm"
  }
}

resource "aws_security_group" "cd_engine_sg" {
  name        = "cd-engine-sg"
  description = "Allow SSH, HTTP, HTTPS"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
