terraform {
  required_providers {
    radosgw = {
      source = "spreadshirt/radosgw"
    }
  }
}

provider "radosgw" {
  endpoint = "http://127.0.0.1:9000"
  # set access_key_id and secret_access_key via ACCESS_KEY_ID and SECRET_ACCESS_KEY env variables
}

variable "user_suffix" {
  type        = string
  description = "suffix to append to generated user"
  default     = ""
}

resource "radosgw_user" "dev_test_user" {
  user_id      = "dev_test${var.user_suffix}"
  display_name = "Ceph dev_test user"
}

resource "radosgw_subuser" "dev_test_subuser_readonly" {
  user_id = "dev_test${var.user_suffix}"
  subuser = "readonly"
  access  = "read"

  depends_on = [radosgw_user.dev_test_user]
}

resource "radosgw_key" "dev_test_default_key" {
  user = "dev_test${var.user_suffix}"

  depends_on = [radosgw_user.dev_test_user]
}

resource "radosgw_key" "dev_test_second_key" {
  user = "dev_test${var.user_suffix}"

  depends_on = [radosgw_user.dev_test_user]
}

resource "radosgw_key" "dev_test_readonly_key" {
  user    = "dev_test${var.user_suffix}"
  subuser = "readonly"

  depends_on = [radosgw_subuser.dev_test_subuser_readonly]
}
