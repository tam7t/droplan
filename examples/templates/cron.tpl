*/5 * * * * root PATH=/sbin DO_KEY=${key} DO_TAG=${tag} /usr/local/bin/droplan >/var/log/droplan.log 2>&1
