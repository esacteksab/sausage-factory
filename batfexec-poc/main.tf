# resource "random_uuid" "uuid" {}
#
# resource "local_file" "uuid" {
#   count    = 3
#   content  = random_uuid.uuid.id
#   filename = "${path.module}/bsTF/${count.index}-uuid.out"
# }
#
# resource "random_pet" "pet" {}
#
# resource "local_file" "pet" {
#   count    = 3
#   content  = random_pet.pet.id
#   filename = "${path.module}/bs/TF/${count.index}-pet.out"
#
# }

# resource "archive_file" "tf_pr" {
#   count       = 10
#   type        = "tar.gz"
#   source_file = "${path.module}/.ctrlc.g0"
#   output_path = "${path.module}/bsTF/${count.index}-ctrlc-g0.tar.gz"
# }
