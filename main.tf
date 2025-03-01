resource "random_uuid" "uuid" {}

resource "local_file" "uuid" {
  count = 1000
  content  = random_uuid.uuid.id
  filename = "${path.module}/bsTF/${count.index}-uuid.out"
}

resource "random_pet" "pet" {}

resource "local_file" "pet" {
  count = 1000
  content  = random_pet.pet.id
  filename = "${path.module}/bs/TF/${count.index}-pet.out"

}

resource "archive_file" "tf_pr" {
  count = 1000
  type        = "tar.gz"
  source_file = "${path.module}/.tf-pr"
  output_path = "${path.module}/bsTF/${count.index}-tf-pr.tar.gz"
}
