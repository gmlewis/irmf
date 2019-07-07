# 001-sphere (ball bearing)

One of the most notoriously-difficult objects to model and create using additive manufacturing
is the perfectly-smooth sphere (*e.g.* a ball bearing, or actually any smooth, curved surface).
STL (a triangle-based representation) is simply the wrong tool for the job.
Not even do voxels solve the problem due to their finite image resolution.

Yet, ironically, the perfect sphere is almost the easiest thing to model using an IRMF shader.

## sphere-1.irmf

Here is an [IRMF shader](sphere-1.irmf) defining a 10mm diameter sphere:

```glsl
/*{
  irmf: "1.0",
  materials: ["AISI 1018 steel"],
  max: [5,5,5],
  min: [-5,-5,-5],
  notes: "Simple IRMF shader - Hello, Sphere!",
  title: "10mm diameter Sphere",
  units: "mm"
}*/

void mainModel4( out vec4 materials, in vec3 xyz ) {
  const float radius = 5.0;  // 10mm diameter sphere.
  float r = length(xyz);  // distance from origin.
  materials[0] = r <= radius ? 1.0 : 0.0; // Only materials[0] is used; the others are ignored.
}
```

* Try loading [sphere-1.irmf](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/examples/001-sphere/sphere-1.irmf) now in the experimental IRMF editor!

## sphere-2.irmf

`sphere-1.irmf` above is fine if your entire model is a sphere, but is not
terribly useful if you would like to make a more complex model out of
one or more spheres. Let's make a `sphere` function that is reusable.

```glsl
/*{
  irmf: "1.0",
  materials: ["AISI 1018 steel"],
  max: [5,5,5],
  min: [-5,-5,-5],
  notes: "Simple IRMF shader - sphere function.",
  title: "10mm diameter Sphere",
  units: "mm"
}*/

float sphere(in vec3 pos, in float radius, in vec3 xyz) {
  xyz -= pos;  // Move sphere into place.
  float r = length(xyz);
  return r <= radius ? 1.0 : 0.0;
}

void mainModel4( out vec4 materials, in vec3 xyz ) {
  const float radius = 5.0;  // 10mm diameter sphere.
  materials[0] = sphere(vec3(), radius, xyz);  // vec3() is [0,0,0] - the origin.
}
```

* Try loading [sphere-2.irmf](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/examples/001-sphere/sphere-2.irmf) now in the experimental IRMF editor!

----------------------------------------------------------------------

# License

Copyright 2019 Glenn M. Lewis. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
