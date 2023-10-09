S3 SSE-C GET
============
This repo is described in the following blog post: *later*

An example application that downloads recursively files under a prefix in
an S3 bucket. The objects are expected to be encrypted with SSE-C. The
key should be given in base64 as an application argument.

Usage:
```
s3-ssec-get mybucket my/prefix AzQwSx468== /tmp/froms3
```


