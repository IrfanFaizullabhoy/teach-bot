# techla-cook

1. Add all required private keys to `/private` folder in the `teach cook` directory.

1. Convenient test command: `ansible all -m ping`

```
ansible-playbook site.yml
```

## Notes

1. Use 2.2 RC1 to address this problem https://github.com/ansible/ansible/issues/17495: `sudo pip install git+https://github.com/ansible/ansible.git@v2.2.0.0-0.1.rc1`

## Todo

1. Get kernel version from system instead of hard-coded variable (docker/tasks/main.yml).

## Resources

1. Ansible best practices: http://docs.ansible.com/ansible/playbooks_best_practices.html

1. Ansible + Docker example: https://github.com/angstwad/docker.ubuntu

1. Ansible + Node.js example: https://github.com/nodesource/ansible-nodejs-role

1. Ansible + Git init bare example: https://github.com/EHER/ansible-git
