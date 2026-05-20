output "instance_ip" {
  description = "Public IP of the AWS VM"
  value       = aws_instance.cd_engine_vm.public_ip
}
