output "instance_id" {
  description = "EC2 instance ID"
  value       = aws_instance.api.id
}

output "instance_public_ip" {
  description = "Public IP for SSH access"
  value       = aws_instance.api.public_ip
}

output "ssh_command" {
  description = "Example SSH command (replace private key path)"
  value       = "ssh -i /path/to/private-key ec2-user@${aws_instance.api.public_ip}"
}
