/*gowdlerize -- execute a program in a sanitized environment

 Copyright (C) 2020 Will Furnass, University of Sheffield

 Redistribution and use in source and binary forms, with or without
 modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice,
    this list of conditions and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
    this list of conditions and the following disclaimer in the documentation
    and/or other materials provided with the distribution.

 3. Neither the name of the copyright holder nor the names of its
    contributors may be used to endorse or promote products derived from this
    software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.

Based on safe_exec C program written by Dave Love, University of Manchester,
which itself was adapted "from the Son of Grid Engine distribution
<https://arc.liv.ac.uk/trac/SGE/> and has, unfortunately, not been properly
reviewed."

This is a wrapper to protect a program run with privileges in a user-influenced
environment by removing security-sensitive environment variables.  (E.g. if you
run something dynamically linked with euid 0 and an LD_LIBRARY_PATH supplied by
the user, all bets are off.)  It works similarly to /usr/bin/busybox env -u ...
(assuming busybox is statically linked) with a -u for all relevant variables
(which may not be very practical).

The hard-wired set of variables below is probably only complete for current
GNU/Linux systems, and maybe not even for all of those.  Extensions for other
systems are welcome.

It is written for use with SGE versions without a fix for the issue configured
to use remote startup daemons (i.e. not using the builtin method in versions
after 6.2) or methods for the prolog etc. when run as a privileged user, e.g.
in sge_conf(5):

 prolog         root@/usr/local/bin/safe_exec /opt/sge/bin/prolog
 ...
 qlogin_daemon  /usr/local/bin/safe_exec /usr/sbin/sshd -i

[For remote startup, you can probably use a static env program to clear the
environment completely, but a prolog, for instance, expects a useful
environment.]

It must be statically linked to be effective, hence this version being written
in Go.
*/
package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/RSE-Sheffield/qsafeexec-rpm/pkg/gowdlerize"
)

func main() {
	// Remove timestamp from Go logger output
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <prog> <arg1> <arg2> ...", os.Args[0])
	}
	prog := os.Args[1]
	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	cmd := exec.Command(prog, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = gowdlerize.CleanEnv(os.Environ())

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
