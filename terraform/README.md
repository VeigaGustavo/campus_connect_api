# Terraform EC2 (SSH Key)

Infra base para subir uma EC2 com acesso por chave SSH (sem login/senha).

## Como usar

1. Copie o arquivo de exemplo:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edite o `terraform.tfvars` e ajuste principalmente:
   - `public_key` com sua chave publica SSH
   - `ssh_allowed_cidr` para seu IP/CIDR
   - `aws_region`, `instance_type` e demais parametros

3. Execute:

```bash
terraform init
terraform plan
terraform apply
```

## Observacoes

- O provisionamento instala Docker na instancia.
- Usuario padrao do Amazon Linux 2023: `ec2-user`.
- A API nao e deployada automaticamente no `user_data`; isso voce pode adaptar depois.
