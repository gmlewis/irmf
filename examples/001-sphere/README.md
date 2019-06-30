# 001-sphere (ball bearing)

One of the most notoriously-difficult objects to model and create using additive manufacturing 
is the perfectly-smooth sphere (*e.g.* a ball bearing, or actually any smooth, curved surface).
STL (a triangle-based representation) is simply the wrong tool for the job.
Not even do voxels solve the problem due to their finite image resolution.

Yet, ironically, the perfect sphere is the very easiest thing to model using an IRMF shader.

Here is an [IRMF shader](sphere.irmf) defining a 10mm diameter sphere:

```glsl
/*{
  author: "Glenn M. Lewis",
  copyright: "Apache-2.0",
  date: "2019-06-30",
  irmf: "1.0",
  materials: ["AISI 1018 steel"],
  max: [5,5,5],
  min: [-5,-5,-5],
  notes: "Simplest-possible IRMF shader - Hello, Sphere!",
  title: "10mm diameter Sphere",
  units: "mm",
  version: "1.0"
}*/

void mainModel4( out vec4 materials, in vec3 xyz ) {
  const float radius = 5.0;  // 10mm diameter sphere.
  float r = length(xyz);  // distance from origin.
  materials[0] = r <= radius ? 1.0 : 0.0; // Only materials[0] is used; the others are ignored.
}
```

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
