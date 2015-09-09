output "client_user_access_key_id" {
  value = "${aws_iam_access_key.client_user_access_key.id}"
}

output "client_user_secret_key" {
  value = "${aws_iam_access_key.client_user_access_key.secret}"
}
