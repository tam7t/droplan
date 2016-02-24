provider "digitalocean" {
  token = "${var.do_token}"
}

resource "digitalocean_ssh_key" "root_ssh" {
  name = "Terraform Root SSH Key"
  public_key = "${file(var.ssh_key_path)}"
}

resource "digitalocean_droplet" "dolan-ubuntu-x64" {
  image = "ubuntu-14-04-x64"
  name = "dolan-ubuntu-x64"
  region = "nyc3"
  size = "512mb"
  private_networking = true
  ssh_keys = ["${digitalocean_ssh_key.root_ssh.id}"]

  provisioner "remote-exec" {
    inline = [
      "apt-get update",
      "apt-get upgrade -y",
      "apt-get install -y unzip"
    ]
  }

  provisioner "remote-exec" {
    inline = [
      "cd /tmp",
      "wget https://github.com/tam7t/dolan/releases/download/v1.0.0/dolan_1.0.0_linux_amd64.zip",
      "unzip dolan_1.0.0_linux_amd64.zip -d /usr/local/bin",
      "rm /tmp/dolan_1.0.0_linux_amd64.zip",
      "echo '${template_file.cron.rendered}' > /etc/cron.d/dolan"
    ]
  }
}

resource "digitalocean_droplet" "dolan-ubuntu-x32" {
  image = "ubuntu-14-04-x32"
  name = "dolan-ubuntu-x32"
  region = "nyc3"
  size = "512mb"
  private_networking = true
  ssh_keys = ["${digitalocean_ssh_key.root_ssh.id}"]

  provisioner "remote-exec" {
    inline = [
      "apt-get update",
      "apt-get upgrade -y",
      "apt-get install -y unzip"
    ]
  }

  provisioner "remote-exec" {
    inline = [
      "cd /tmp",
      "wget https://github.com/tam7t/dolan/releases/download/v1.0.0/dolan_1.0.0_linux_386.zip",
      "unzip dolan_1.0.0_linux_386.zip -d /usr/local/bin",
      "rm /tmp/dolan_1.0.0_linux_386.zip",
      "echo '${template_file.cron.rendered}' > /etc/cron.d/dolan"
    ]
  }
}

resource "digitalocean_droplet" "dolan-fedora-x64" {
  image = "fedora-23-x64"
  name = "dolan-fedora-x64"
  region = "nyc3"
  size = "512mb"
  private_networking = true
  ssh_keys = ["${digitalocean_ssh_key.root_ssh.id}"]

  provisioner "remote-exec" {
    inline = [
      "dnf -y install unzip wget"
    ]
  }

  provisioner "remote-exec" {
    inline = [
      "curl -O -L https://github.com/tam7t/dolan/releases/download/v1.0.0/dolan_1.0.0_linux_amd64.zip",
      "unzip dolan_1.0.0_linux_amd64.zip -d /usr/local/bin",
      "rm dolan_1.0.0_linux_amd64.zip",
      "echo '${template_file.cron.rendered}' > /etc/cron.d/dolan"
    ]
  }
}

resource "template_file" "cron" {
  template = "${file("${path.module}/templates/cron.tpl")}"

  vars {
    key = "${var.dolan_token}"
  }
}
