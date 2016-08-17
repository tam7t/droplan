#cloud-config

coreos:
  units:
  - name: droplan-setup.service
    command: start
    content: |
      [Unit]
      Description=setup droplan iptable rules for docker

      [Service]
      Type=oneshot
      After=docker.service
      ExecStart=/usr/bin/sh -c "docker ps; \
        iptables -N droplan-peers; \
        iptables -I FORWARD 1 -i eth1 -j DROP; \
        iptables -I FORWARD 1 -i eth1 -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT; \
        iptables -I FORWARD 1 -j droplan-peers"
  - name: droplan.service
    command: start
    content: |
      [Unit]
      Description=updates iptables with peer droplets
      After=droplan-setup.service
      Requires=docker.service

      [Service]
      Type=oneshot
      Environment=DO_KEY=${key}
      Environment=DO_TAG=${tag}
      ExecStart=/usr/bin/docker run --rm --net=host --cap-add=NET_ADMIN -e DO_KEY -e DO_TAG tam7t/droplan:latest
  - name: droplan.timer
    command: start
    content: |
      [Unit]
      Description=Run droplan.service every 5 minutes

      [Timer]
      OnCalendar=*:0/5
