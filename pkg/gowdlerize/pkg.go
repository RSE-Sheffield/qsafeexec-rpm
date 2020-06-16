package gowdlerize

import (
	"regexp"
	"runtime"
	"strings"
)

var stringSet = map[string]struct{}{}

// CleanEnv sanitises an environment.
func CleanEnv(env []string) []string {
	var out []string

	for _, e := range env {
		pair := strings.SplitN(e, "=", 2)
		if isUnsafe(pair[0]) {
			continue
		}
		if isBadLocale(pair[0], pair[1]) {
			continue
		}
		out = append(out, e)
	}
	return out
}

func isUnsafe(name string) bool {
	switch name {
	case
		/* These variables are regarded as unsafe in suid programs by
		   GNU libc (sysdeps/generic/unsecvars.h as of v2.31). */
		"GCONV_PATH",
		"GETCONF_DIR",
		"GLIBC_TUNABLES",
		"HOSTALIASES",
		"LD_AUDIT",
		"LD_DEBUG",
		"LD_DEBUG_OUTPUT",
		"LD_DYNAMIC_WEAK",
		"LD_HWCAP_MASK",
		"LD_LIBRARY_PATH",
		"LD_ORIGIN_PATH",
		"LD_PRELOAD",
		"LD_PROFILE",
		"LD_SHOW_AUXV",
		"LD_USE_LOAD_BIAS",
		"LOCALDOMAIN",
		"LOCPATH",
		"MALLOC_TRACE",
		"NIS_PATH",
		"NLSPATH",
		"RESOLV_HOST_CONF",
		"RES_OPTIONS",
		"TMPDIR",
		"TZDIR",
		/* Fixme:  Check if IFS actually needs to be set on relevant systems
		   <http://www.dwheeler.com/secure-programs/Secure-Programs-HOWTO/environment-variables.html#ENV-VAR-SOLUTION>.  */
		"IFS",
		/* From sudo (plugins/sudoers/env.c as of v1.9.1): */
		"CDPATH",
		"PATH_LOCALE",
		"SHLIB_PATH", // on HPUX
		"LIBPATH",    // on AIX
		"AUTHSTATE",  // on AIX
		"KRB5_CONFIG*",
		"KRB5_KTNAME",
		"VAR_ACE",
		"USR_ACE",
		"DLC_ACE",
		"TERMINFO",          /* terminfo, exclusive path to terminfo files */
		"TERMINFO_DIRS",     /* terminfo, path(s) to terminfo files */
		"TERMPATH",          /* termcap, path(s) to termcap files */
		"TERMCAP",           /* XXX - only if it starts with '/' */
		"ENV",               /* ksh, file to source before script runs */
		"BASH_ENV",          /* bash, file to source before script runs */
		"PS4",               /* bash, prefix for lines in xtrace mode */
		"GLOBIGNORE",        /* bash, globbing patterns to ignore */
		"BASHOPTS",          /* bash, initial "shopt -s" options */
		"SHELLOPTS",         /* bash, initial "set -o" options */
		"JAVA_TOOL_OPTIONS", /* java, extra command line options */
		"PERLIO_DEBUG",      /* perl, debugging output file */
		"PERLLIB",           /* perl, search path for modules/includes */
		"PERL5LIB",          /* perl 5, search path for modules/includes */
		"PERL5OPT",          /* perl 5, extra command line options */
		"PERL5DB",           /* perl 5, command used to load debugger */
		"FPATH",             /* ksh, search path for functions */
		"NULLCMD",           /* zsh, command for null file redirection */
		"READNULLCMD",       /* zsh, command for null file redirection */
		"ZDOTDIR",           /* zsh, search path for dot files */
		"TMPPREFIX",         /* zsh, prefix for temporary files */
		"PYTHONHOME",        /* python, module search path */
		"PYTHONPATH",        /* python, search path */
		"PYTHONINSPECT",     /* python, allow inspection */
		"PYTHONUSERBASE",    /* python, per user site-packages directory */
		"RUBYLIB",           /* ruby, library load path */
		"RUBYOPT":           /* ruby, extra command line options */
		return true
	}

	// "*=()*",			/* bash functions */  /* TODO */

	/* Check for "LD_AUDIT", "LD_DEBUG", "LD_DEBUG_OUTPUT",
	   "LD_DYNAMIC_WEAK", "LD_LIBRARY_PATH", "LD_ORIGIN_PATH",
	   "LD_PRELOAD", "LD_PROFILE", "LD_SHOW_AUXV", "LD_USE_LOAD_BIAS",
	   plus others not from glibc e.g.
		LD_LIBRARY_PATH &c, DYLD_LIBRARY_PATH, _RLD_*, LDR_* */
	if strings.HasPrefix(name, "LD_") {
		return true
	}
	if strings.HasPrefix(name, "_RLD_") {
		return true
	}
	if runtime.GOOS == "darwin" && strings.HasPrefix(name, "DYLD_") {
		/* DYLD_{FALLBACK_,}FRAMEWORK_PATH,
		   {FALLBACK_,}DYLD_LIBRARY_PATH seem to be relevant.  */
		return true
	}
	// On AIX
	if name == "LIBPATH" || strings.HasPrefix(name, "LDR_") {
		return true
	}
	// On HPUX
	if name == "SHLIB_PATH" || strings.HasPrefix(name, "_HP_DLD") {
		/* _HP_DLDOPTS is loader options; guess a prefix.  */
		return true
	}

	return false
}

func isBadLocale(name string, val string) bool {
	validLocale := regexp.MustCompile(`^[A-Za-z_]+=[A-Za-z][-A-Za-z0-9_,+@.=]*$`)

	if !(name == "LANGUAGE" || name == "LANG" || name == "LINGUAS" || strings.HasPrefix(name, "LC_")) {
		return false
	}

	if validLocale.MatchString(val) {
		return false
	}
	return true
}
