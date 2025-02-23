resource "random_uuid" "uuid" {}

resource "local_file" "uuid" {
  content  = random_uuid.uuid.id
  filename = "${path.module}/uuid.out"
}

resource "random_pet" "pet" {}

resource "local_file" "pet" {
  content  = random_pet.pet.id
  filename = "${path.module}/pet.out"

}

resource "archive_file" "tf_pr" {
  type        = "tar.gz"
  source_file = "${path.module}/.tf-pr"
  output_path = "${path.module}/tf-pr.tar.gz"
}
