This example shows how to create a VM and install [yunohost](https://yunohost.org/).
This include:
- client-side generated keypair
- security group and rules (which are open because yunohost manage its own firewalling rules locally.
- Public IP
- Resized bootdisk volume

Once created, you can connect to your new yunohost instance throw its web interface and finalize installation.
See more details on [yunohost documentation](https://yunohost.org/fr/install/hardware:vps_debian).
Note that sshd configuration is overwritten but you will be still able to connect with `admin` user with the password set on installation finalization.
