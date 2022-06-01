# Orbot IPtProxy

This extends IPtProxy that includes additional functionality (go-tun2socks) needed for Orbot specifically.

This is needed due to the fact that an Android app may only have one bundled gomobile built AAR at a time. This repository makes it trivial to integrate new releases of IPtProxy with additional go code that Orbot needs. 

To update IPtProxy releases (obfs4 and snowflake) update the submodule to point at a given release. Right now this is set to IPtProxy 1.6.0

Additional go code for Orbot can be added to the `./OrbotTun.go` directory. 

To build an `AAR` and sources `JAR` run `./build-orbot.sh` after setting the `ANDROID_HOME` and `ANDROID_NDK_HOME` environment variables.
