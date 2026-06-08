# provider example

Can be run by starting `docker compose up` to get a local Ceph setup and then running
`terraform apply`:

```
ACCESS_KEY_ID=RMkni81ukvCYTLCjk62d SECRET_ACCESS_KEY=k8xeC8Kb62PMSXglkeuS6kLLjOHRp6y5LMntsUAR TF_CLI_CONFIG_FILE=local-dev.tfrc terraform apply
```

For repeated runs, there are several approaches:

- always clean up (remove `terraform.tfstate*` and created `dev_test` user)
- or run with the following command to always generate a new user and use a fresh Terraform state:

    ```
    ACCESS_KEY_ID=RMkni81ukvCYTLCjk62d SECRET_ACCESS_KEY=k8xeC8Kb62PMSXglkeuS6kLLjOHRp6y5LMntsUAR TF_CLI_CONFIG_FILE=local-dev.tfrc terraform apply -state=/tmp/rgw-tf$RANDOM.tfstate -auto-approve -var user_suffix=$RANDOM
    ```
