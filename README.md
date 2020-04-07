# RPM packaging of Son of Grid Engine's safe_exec wrapper

Renamed to `qsafeexec` to more clearly indicate its origins (Grid Engine commands all start with `q`).

See `qsafeexec.1` man page for more info.

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
