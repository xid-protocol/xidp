# xid-protocol
Enterprise-grade, Infinitely-Scalable, Distributed Identity-Data Protocol

# Install
Prepare a default configuration, uncomment add configure as needed:
```
mkdir -p ~/.config/xidp && cat <<'EOF' > ~/.config/xidp/config.yml

Server:
  port: 9527

#Feilian:
#  access_key_id: xxxx
#  access_key_secret: xxxx
#  endpoint: http://

#Jumpserver:
#  endpoint: https://jumpserver.xx.com
#  access_key_id: xxx
#  access_key_secret: xxx

Log:
  path: /var/log/xidp/xidp.log
  max_size: 10
  max_age: 30
  max_backups: 3

MongoDB:
  uri: mongodb://username:password@192.0.8.11:27017/?authSource=admin
  database: xid-protocol
EOF
```

custom config path

```
./xidp -c config.yaml
```