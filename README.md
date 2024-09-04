# cilium-bpf-hack

This is a bunch of scripts to visualize the Cilium BPF codebase. 

WARNING: The software here is a hack, so be aware of its general instability and
hackery.

## Usage

The general flow is:

* Generate an LLVM bitcode output from the BPF code.
* Generate an associated annotations file from the source code using `cmd/genam/genam.sh`.
* Generate the graphviz PDF output using `cmd/cfg`.

### Generating the BPF LLVM object code

WARNING: total hack here.

Make the bpf test object file. Right now the way I do this is to grab the
command line for the fully expanded options from running `make` and change the
optimization to `-O=0 -emit-llvm` and remove inlining:

```
$ clang -O0 -DSKIP_DEBUG=1 \
  -DENABLE_IPV4=1 -DENABLE_IPV6=1 -DENABLE_ROUTING=1 ... \
  -I/home/bowei/work/cilium/bpf/include \
  -I/home/bowei/work/cilium/bpf \
  -D__NR_CPUS__=8 -g --target=bpf -std=gnu89 -nostdinc \
  -Wall -W... \
  -mcpu=v3 -emit-llvm -S -c bpf_lxc.c \
  -o bpf_lxc.ll
```

Inline removal patch:

```
diff --git a/bpf/include/bpf/compiler.h b/bpf/include/bpf/compiler.h
index d685e454e8..2b0c4cc673 100644
--- a/bpf/include/bpf/compiler.h
+++ b/bpf/include/bpf/compiler.h
@@ -42,7 +42,9 @@
 #endif
 
 #undef __always_inline         /* stddef.h defines its own */
-#define __always_inline                inline __attribute__((always_inline))
+//#define __always_inline              inline __attribute__((always_inline))
+
+#define __always_inline
 
 #ifndef __stringify
 # define __stringify(X)                #X
diff --git a/bpf/lib/source_info.h b/bpf/lib/source_info.h
index 5cdc083f4e..79324a2d18 100644
--- a/bpf/lib/source_info.h
+++ b/bpf/lib/source_info.h
@@ -3,7 +3,7 @@
 #pragma once
 
 #ifndef BPF_TEST
-#define __MAGIC_FILE__ (__u8)__id_for_file(__FILE_NAME__)
+#define __MAGIC_FILE__ (__u8)__id_for_file(__FILE__)
 #define __MAGIC_LINE__ __LINE__
 #else
 /* bpf tests assert that metrics get updated by performing a map lookup.
```

This will generate a `bpf_lxc.ll` that has the object file as LLVM bitcode.

### Generating diagrams

#### Source annotations

`genam` generates an annotations file from the source code. This is a file that
contains a mapping of `file:line` with comments/annotations in the
resulting diagrams.

```
$ bash cmd/genan/genan.sh ../cilium/bpf > annotations.txt
```

#### Control flow output

```
# $1 is the name of the function to start the diagram from. E.g. cil_from_container.
$ ./cfg -mode rawcg \
  -in bpf_lxc.ll \
  -an annotations1.txt \
  -an annotations2.txt \
  -start $1 \
  > /tmp/out.gv

$ dot -Tpdf /tmp/out.gv -o out.pdf
```