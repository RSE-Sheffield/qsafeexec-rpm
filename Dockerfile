FROM centos:7
ENV rhelmajor 7

# RPM name (need a spec called $name.spec in the working directory) 
ARG name
# RPM version
ARG version
# RPM release
ARG release=1

# Check we've been passed the necessary arguments to build a RPM at build-time:
RUN test -n "$name" && test -n "$version"

# Our inputs: RPM spec and custom manpage
ENV specfilepath=/home/unpriv/rpmbuild/SPECS/$name.spec \
    manfilepath=/home/unpriv/rpmbuild/SOURCES/$name.1 \
    tarballname=safe_exec.tar.gz
ENV tarballurl https://arc.liv.ac.uk/downloads/SGE/support/$tarballname
# Our output
ENV rpmfilepath /home/unpriv/rpmbuild/RPMS/x86_64/$name-$version-$release.el$rhelmajor.x86_64.rpm

# Install packages needed to build RPMs
RUN yum update -y && \
    yum install -y \
        rpm-build \
        rpmrebuild \
        yum-utils \
        rpmdevtools 

# Set up the RPM build environment as an unpriv user
RUN useradd -ms /bin/bash unpriv
USER unpriv
WORKDIR /home/unpriv
RUN rpmdev-setuptree
COPY --chown=unpriv $name.spec $specfilepath
COPY --chown=unpriv $name.1 $manfilepath
ADD --chown=unpriv $tarballurl /home/unpriv/rpmbuild/SOURCES/$tarballname

# Install build depencencies
USER root
RUN yum-builddep -y $specfilepath

# Build the RPM 
USER unpriv
RUN rpmbuild -bb $specfilepath

# Checks: 
# - we can install, uninstall and reinstall the RPM
# - the qsafeexec executable is statically linked
# - the qsafeexec executable scrubs undesirable vars from the wrapped executable's env
USER root
RUN rpm -ivh $rpmfilepath && \
    rpm -e $name && \
    rpm -ivh $rpmfilepath && \
    file /opt/qsafeexec/bin/qsafeexec | grep -q 'statically linked' && \
    export LD_PRELOAD=/usr/lib64/libxml2.so.2 && /opt/qsafeexec/bin/qsafeexec sh -c 'echo $LD_PRELOAD' | grep -vq libxml2 && \
    export SOMEVAR=/usr/lib64/libxml2.so.2 && /opt/qsafeexec/bin/qsafeexec sh -c 'echo $SOMEVAR' | grep -q libxml2
