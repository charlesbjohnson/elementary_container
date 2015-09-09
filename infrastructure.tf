provider "aws" {
    access_key = "${var.aws_access_key_id}"
    secret_key = "${var.aws_secret_key}"
    region = "${var.aws_region}"
}

resource "aws_iam_user" "client_user" {
    name = "cbjohnson-os-client"
    path = "/os/"
}

resource "aws_iam_access_key" "client_user_access_key" {
    user = "${aws_iam_user.client_user.name}"
}

resource "aws_iam_role" "client_role" {
    name = "os-client-role"
    path = "/os/"
    assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Principal": {
        "AWS": "arn:aws:iam::${var.aws_account_id}:root"
      }
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "client_role_policy" {
    name = "os-client-role-image-bucket-policy"
    role = "${aws_iam_role.client_role.id}"
    policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:ListBucket",
        "s3:PutObject"
      ],
      "Resource": [
        "arn:aws:s3:::${var.aws_os_images_bucket}",
        "arn:aws:s3:::${var.aws_os_images_bucket}/*"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_user_policy" "client_user_policy" {
    name = "os-client-user-policy"
    user = "${aws_iam_user.client_user.name}"
    policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Resource": "${aws_iam_role.client_role.arn}"
    }
  ]
}
EOF
}

resource "aws_s3_bucket" "os_image_bucket" {
  bucket = "${var.aws_os_images_bucket}"
}
