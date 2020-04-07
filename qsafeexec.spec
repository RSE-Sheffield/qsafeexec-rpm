Name:           qsafeexec
Version:        0.1
Release:        1%{?dist}
Summary:        Statically-linked wrapper for dynamicly linked executables that must be run with elevated privileges.
License:        GPLv3+
Packager:       Will Furnass <w.furnass@sheffield.ac.uk>
URL:            https://arc.liv.ac.uk/SGE/htmlman/htmlman5/sge_conf.html
Source0:        safe_exec.tar.gz
%define         SHA256SUM0 3246439cc8147ce99bbb5df060fbc28fbd867fb0c640ea228b29d68c98b83c3a

BuildRequires:  coreutils make gcc glibc-static
#Requires:       NONE       

%description
The qsafeexec wrapper program is intended to address a gridengine
security issue (CVE-2012-0208), but might be of more general use.
Guards against LD_PRELOAD/LD_LIBRARY_PATH attacks and certain other attacks.

%prep
echo "%SHA256SUM0  %SOURCE0" | sha256sum -c -

%setup -q -n safe_exec

%build
make safe_exec %{?_smp_mflags}

%install
rm -rf $RPM_BUILD_ROOT

mkdir -p $RPM_BUILD_ROOT/opt/%{name}/bin
install -m 0755 safe_exec $RPM_BUILD_ROOT/opt/%{name}/bin/%{name}

install -m 0644 README $RPM_BUILD_ROOT/opt/%{name}/

mkdir -p $RPM_BUILD_ROOT/opt/%{name}/man/man1
install -m 0644 %{_sourcedir}/%{name}.1 $RPM_BUILD_ROOT/opt/%{name}/man/man1/%{name}.1 

%clean
rm -rf $RPM_BUILD_ROOT

%files
/opt/%{name}/bin/%{name}
%doc /opt/%{name}/README
%doc /opt/%{name}/man/man1/%{name}.1 

%changelog

* Mon Apr 06 2020 Will Furnass <w.furnass@sheffield.ac.uk> 0.1-1
- Package safe_exec from https://arc.liv.ac.uk/downloads/SGE/support/ (date in README: 2012-04)
