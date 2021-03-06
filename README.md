# dockviz

Visualizing Docker Data

This command takes the raw Docker JSON and visualizes it in various ways.

For image information, output can be formatted as
[Graphviz](http://www.graphviz.org) or as a tree in the terminal.

For container information, only Graphviz output has been implemented.

# Examples

## Containers

Currently, containers are visualized with labeled lines for links.  Containers that aren't running are greyed out.

![](sample/containers.png "Container")

## Images

Image info is visualized with lines indicating parent images:

![](sample/images.png "Image")

Or as a tree in the terminal:

```
└─511136ea3c5a Virtual Size: 0.0 B
  |─f10ebce2c0e1 Virtual Size: 103.7 MB
  | └─82cdea7ab5b5 Virtual Size: 103.9 MB
  |   └─5dbd9cb5a02f Virtual Size: 103.9 MB
  |     └─74fe38d11401 Virtual Size: 209.6 MB Tags: ubuntu:12.04, ubuntu:precise
  |─ef519c9ee91a Virtual Size: 100.9 MB
  | └─07302703becc Virtual Size: 101.2 MB
  |   └─cf8dc907452c Virtual Size: 101.2 MB
  |     └─a7cf8ae4e998 Virtual Size: 171.3 MB Tags: ubuntu:12.10, ubuntu:quantal
  |       |─e18d8001204e Virtual Size: 171.3 MB
  |       | └─d0525208a46c Virtual Size: 171.3 MB
  |       |   └─59dac4bae93b Virtual Size: 242.5 MB
  |       |     └─89541b3b35f2 Virtual Size: 511.8 MB
  |       |       └─7dac4e98548e Virtual Size: 511.8 MB
  |       |         └─341d0cc3fac8 Virtual Size: 511.8 MB
  |       |           └─2f96171d2098 Virtual Size: 511.8 MB
  |       |             └─67b8b7262a67 Virtual Size: 513.7 MB
  |       |               └─0fe9a2bc50fe Virtual Size: 513.7 MB
  |       |                 └─8c32832f07ba Virtual Size: 513.7 MB
  |       |                   └─cc4e1358bc80 Virtual Size: 513.7 MB
  |       |                     └─5c0d04fba9df Virtual Size: 513.7 MB Tags: nate/mongodb:latest
  |       └─398d592f2009 Virtual Size: 242.2 MB
  |         └─0cd8e7f50270 Virtual Size: 243.6 MB
  |           └─594b6f8e6f92 Virtual Size: 243.6 MB
  |             └─f832a63e87a4 Virtual Size: 243.6 MB Tags: redis:latest
  └─02dae1c13f51 Virtual Size: 98.3 MB
    └─e7206bfc66aa Virtual Size: 98.5 MB
      └─cb12405ee8fa Virtual Size: 98.5 MB
        └─316b678ddf48 Virtual Size: 169.4 MB Tags: ubuntu:13.04, ubuntu:raring
```

# Running

Currently, this only works when the remote API is listening on TCP.  Soon, the Docker command line will allow dumping the image JSON.

```
$ curl -s http://localhost:4243/images/json?all=1 | ./dockviz images --dot | dot -Tpng -o images.png
```

# Download

For now, download binaries from Gobuild: <http://gobuild.io/download/github.com/justone/dockviz>
