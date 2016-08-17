provider "digitalocean" {
  token = "${var.do_token}"
}

resource "digitalocean_tag" "droplan" {
  name = "droplantag"
}

resource "digitalocean_ssh_key" "root_ssh" {
  name = "Terraform Root SSH Key"
  public_key = "${file(var.ssh_key_path)}"
}

resource "digitalocean_droplet" "droplan-ubuntu-x64" {
  image = "ubuntu-14-04-x64"
  name = "droplan-ubuntu-x64"
  region = "nyc3"
  size = "512mb"
  private_networking = true
  ssh_keys = ["${digitalocean_ssh_key.root_ssh.id}"]
  tags = ["${digitalocean_tag.droplan.id}"]

  provisioner "remote-exec" {
    inline = [
      "cd /tmp",
      "curl -O -L https://github.com/tam7t/droplan/releases/download/v1.1.0/droplan_1.1.0_linux_amd64.tar.gz",
      "tar -zxf droplan_1.1.0_linux_amd64.tar.gz -C /usr/local/bin",
      "rm /tmp/droplan_1.1.0_linux_amd64.tar.gz",
      "echo '${data.template_file.cron.rendered}' > /etc/cron.d/droplan"
    ]
  }
}

resource "digitalocean_droplet" "droplan-ubuntu-x32" {
  image = "ubuntu-14-04-x32"
  name = "droplan-ubuntu-x32"
  region = "nyc3"
  size = "512mb"
  private_networking = true
  ssh_keys = ["${digitalocean_ssh_key.root_ssh.id}"]
  tags = ["${digitalocean_tag.droplan.id}"]

  provisioner "remote-exec" {
    inline = [
      "cd /tmp",
      "curl -O -L https://github.com/tam7t/droplan/releases/download/v1.1.0/droplan_1.1.0_linux_386.tar.gz",
      "tar -zxf droplan_1.1.0_linux_386.tar.gz -C /usr/local/bin",
      "rm /tmp/droplan_1.1.0_linux_386.tar.gz",
      "echo '${data.template_file.cron.rendered}' > /etc/cron.d/droplan"
    ]
  }
}

resource "digitalocean_droplet" "droplan-fedora-x64" {
  image = "fedora-23-x64"
  name = "droplan-fedora-x64"
  region = "nyc3"
  size = "512mb"
  private_networking = true
  ssh_keys = ["${digitalocean_ssh_key.root_ssh.id}"]
  tags = ["${digitalocean_tag.droplan.id}"]

  provisioner "remote-exec" {
    inline = [
      "curl -O -L https://github.com/tam7t/droplan/releases/download/v1.1.0/droplan_1.1.0_linux_amd64.tar.gz",
      "tar -zxf droplan_1.1.0_linux_amd64.tar.gz -C /usr/local/bin",
      "rm droplan_1.1.0_linux_amd64.tar.gz",
      "echo '${data.template_file.cron.rendered}' > /etc/cron.d/droplan"
    ]
  }
}

resource "digitalocean_droplet" "droplan-coreos" {
  image = "coreos-stable"
  name = "droplan-coreos"
  region = "nyc3"
  size = "512mb"
  private_networking = true
  ssh_keys = ["${digitalocean_ssh_key.root_ssh.id}"]
  tags = ["${digitalocean_tag.droplan.id}"]

  user_data = "${data.template_file.cloudinit.rendered}"
}

resource "digitalocean_droplet" "droplan-ubuntu-x64-notag" {
  image = "ubuntu-14-04-x64"
  name = "droplan-ubuntu-x64-notag"
  region = "nyc3"
  size = "512mb"
  private_networking = true
  ssh_keys = ["${digitalocean_ssh_key.root_ssh.id}"]

  provisioner "remote-exec" {
    inline = [
      "cd /tmp",
      "curl -O -L https://github.com/tam7t/droplan/releases/download/v1.1.0/droplan_1.1.0_linux_amd64.tar.gz",
      "tar -zxf droplan_1.1.0_linux_amd64.tar.gz -C /usr/local/bin",
      "rm /tmp/droplan_1.1.0_linux_amd64.tar.gz",
      "echo '${data.template_file.cron.rendered}' > /etc/cron.d/droplan"
    ]
  }
}

data "template_file" "cloudinit" {
  template = "${file("${path.module}/templates/cloudinit.tpl")}"

  vars {
    key = "${var.droplan_token}"
    tag = "${digitalocean_tag.droplan.id}"
  }
}

data "template_file" "cron" {
  template = "${file("${path.module}/templates/cron.tpl")}"

  vars {
    key = "${var.droplan_token}"
    tag = "${digitalocean_tag.droplan.id}"
  }
}
