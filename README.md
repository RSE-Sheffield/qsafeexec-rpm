# RPM packaging of Son of Grid Engine's safe_exec wrapper

`qsafeexec`: statically-linked wrapper for dynamicly linked executables that must be run with elevated privileges.
Renamed to `qsafeexec` to more clearly indicate its origins (Grid Engine commands all start with `q`).

`safe_exec` is copyright (C) 2011 Dave Love, University of Liverpool (`d.love@liverpool.ac.uk`).

> The `qsafeexec` wrapper program is intended to address a Grid Engine
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

`qsafeexec` cleans the environment for the wrapped executable by
scrubbing several environment variables: `LD_*`, `DYLD_LIBRARY_PATH`,
`_RLD_*`, `LDR_*`, `LOCALDOMAIN`, `LOCPATH`, `MALLOC_TRACE`, `NIS_PATH`,
`NLSPATH`, `RESOLV_HOST_CONF`, `RES_OPTIONS`, `TMPDIR`, `TZDIR`, `IFS`,
`SHELLOPTS`, `PS4`, `PATH_LOCALE`, `PATH_LOCALE`, `TERMINFO`, `TERMINFO_DIRS`,
`TERMPATH`, `TERMCAP`, `ENV`, `BASH_ENV`, `KRB5_CONFIG`, `KRB5_KTNAME`,
`JAVA_TOOL_OPTIONS` (plus `AUTHSTATE` on AIX), and by performing sanity
checks on locale-related variables.

> The code is adapted from the [Son of Grid Engine distribution](https://arc.liv.ac.uk/trac/SGE/)
and has, unfortunately, not been properly reviewed. Please report any problems directly to the author or
via the issue tracker at the site above.

## RPM structure

The `.spec` file in this repo will create an RPM that installs the qsafeexec binary and man page under `/opt/qsafeexec`.

## Building RPM for Centos 7

```sh
docker build  --build-arg=name=qsafeexec --build-arg=version=0.1 --build-arg=release=1 .
docker create $image  # returns container ID
docker cp $container_id:/home/unpriv/rpmbuild/RPMS/x86_64/qsafeexec-0.1-1.el7.x86_64.rpm .
docker rm $container_id
```

## Building RPM for Centos 6

* Modify the RHEL version in the first two lines of the Dockerfile
* Repeat the Centos 7 steps shown above but replace `el7` with `el6`
