# go-cstring

C strings whose memory is managed by the Go runtime.

There are two primary advantages to this library over using C.CString:

1. An extra byte is allocated for the null-terminator byte, so converting a
   CString to a Go string representation or C string representation can be done
   in constant time without additional memory allocations.
1. A CString is just a regular Go slice under the hood, so the memory is
   automatically freed when the CString is no longer in use.
