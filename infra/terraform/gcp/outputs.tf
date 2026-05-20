output "instance_ip" {
  description = "Public IP of the GCP VM"
  value       = google_compute_instance.cd_engine_vm.network_interface[0].access_config[0].nat_ip
}
