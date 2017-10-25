# grplog
package grplog creates a grouped set of convenience hooks into stdlib
log.logger for various levels of logging: Trace, Debug, Info, Notice,
Warning, Alert, Error, Critical, Emergency. Having them grouped easily
allows multiple sets within a program as well as providing a way to easily
grep the output coming from that group. i.e. glog:TRACE: glog:DEBUG versus
blog:TRACE: blog:ERROR
