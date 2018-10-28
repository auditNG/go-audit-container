#!/bin/bash
set -e

cp ../go-audit-container go_audit_container_debian/usr/local/bin/
mv go_audit_container_debian/usr/local/bin/go-audit-container go_audit_container_debian/usr/local/bin/go-audit
cp ../go-audit.yaml go_audit_container_debian/etc/
cp ../examples/go-audit/systemd.service go_audit_container_debian/etc/systemd/system/goaudit.service
dpkg-deb --build go_audit_container_debian
