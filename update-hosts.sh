#!/bin/bash
## Update /etc/hosts with the contents of /etc/hosts.d/*
cat /etc/hosts.d/* > /etc/hosts
