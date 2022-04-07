# Mat

This is an experimental use of Go generics to implementation matrix and vector
operations over elements from a "Field".

All of this is subject to change at any moment.

TODO: would be interesting to see if operations could be "thunked",
meaning evaluation deferred until it's useful to do it, in hopes
of doing reassociation and temporary layout choices that produce
faster results.
