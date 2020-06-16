# gowdlerize: sanitise environments when running binaries with elevated privileges and user-supplied environments

This is a statically-linked wrapper for dynamically linked executables that must be run with elevated privileges.

Based on `safe_exec` by Dave Love, University of Liverpool (`d.love@liverpool.ac.uk`), who said about that:

> The `safe_exec` wrapper program is intended to address a Grid Engine
> security issue (CVE-2012-0208), but might be of more general use.
> 
> The issue concerns a trivial remote root on execution hosts for
> unprivileged valid users on submit hosts, but only for non-default
> configurations in recent versions. It also affects versions before SGE
> 6.2 for probably any configuration with interactive queues. The wrapper
> is intended for users who can't install a Grid Engine version with a
> fix (such as Son of Grid Engine 8.0.0e) and who need to run external
> remote startup programs (like `ssh`) or privileged Grid Engine prolog
> scripts etc. It implements the approach of Son of Grid Engine 8.0.0e,
> but as an add-on which executes the target command in a sanitized
> environment. For more information about the issue, see the security
> notes in `remote_startup` (5) and `sge_conf` (5)
> 
> It should be used to protect Grid Engine remote startup programs (when
> not using the `builtin` method), and also the `prolog` etc. methods if they
> are run as a privileged user (the `<user\>@` prefix). Configure it
> something like this with:
> 
>     qconf -mconf
> 
> or `qmon`(1)'s *Cluster Configuration* (see `sge_conf` (5)):
> 
>     ...
>     prolog root@/opt/qsafeexec/bin/qsafeexec /opt/sge/bin/prolog
>     ...
>     qlogin_daemon /opt/qsafeexec/bin/qsafeexec /usr/sbin/sshd -i
>     ...
> 
> and similarly for other methods you need.

`gowderlize` cleans the environment for the wrapped executable by
scrubbing several environment variables (inc those deemed bad
by glibc for supplying to setuid programs and those blacklisted by sudo)
and by performing sanity checks on locale-related variables.


## Building

```sh
go build -o gowdlerize cmd/gowdlerize/main.go
```

## Building

```sh
go test ./...
```
