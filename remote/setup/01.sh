#!/bin/bash
set -eu

## VARIABLES

# time zone
TIMEZONE=America/New_York

# username
USERNAME=greenlight

# get database password
read -p "Enter password for greenlight DB user:" DB_PASSWORD

# set locale to english
export LC_ALL=en_US.UTF-8

## SCRIPT

# enable the universe repo
add-apt-repository --yes universe

# update all software packages
apt update
apt --yes -o Dpkg::Options::="--force-confnew" upgrade

# set system timezone and install locales
timedatectl set-timezone ${TIMEZONE}
apt --yes install locales-all

# add new user with sudo privileges
useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"

# force password reset
passwd --delete "${USERNAME}"
chage --lastday 0 "${USERNAME}"

# copy ssh keys from root to new user
rsync --archive --chown=${USERNAME}:${USERNAME} /root/.ssh /home/${USERNAME}

# configure firewall to allow SSH, HTTP and HTTPS traffic
ufw allow 22
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

# install fail2ban
apt --yes install fail2ban
apt --yes install curl

# install migrate tools
curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$os-$arch.tar.gz | tar xvzmv migrate.linux-amd64 /usr/local/bin/migrate
mv migrate.linux-amd64 /usr/local/bin/migrate

# install postgres
apt --yes install postgresql

# set up greenlight DB and create user with the password entered earlier
sudo -i -u postgres psql -c "CREATE DATABASE greenlight"
sudo -i -u postgres psql -d greenlight -c "CREATE EXTENSION IF NOT EXISTS citext"
sudo -i -u postgres psql -d greenlight -c "CREATE ROLE greenlight WITH LOGIN PASSWORD '${DB_PASSWORD}'"

# add database dsn to environment variable
echo "GREENLIGHT_DB_DSN='postgres://greenlight:${DB_PASSWORD}@localhost/greenlight?sslmode=disable" >> /etc/environment

# install caddy
apt --yes install -y debian-keyring debian-archive-keyring apt-transport-https
curl -L https://dl.cloudsmith.io/public/caddy/stable/gpg.key | sudo apt-key add -
curl -L https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt | sudo tee -a /etc/apt/sources.list.d/caddy-stable.list
apt update
apt --yes install caddy

echo "Script complete! Rebooting"
reboot