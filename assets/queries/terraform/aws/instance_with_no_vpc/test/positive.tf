resource "aws_instance" "example" {
  ami = "ami-003634241a8fcdec0"

  instance_type = "t2.micro"

}