# 003-coil-square-face

Another surprisingly-simple model is a helical coil with a square cross-section
face.

## coil-1.irmf

```glsl
/*{
  irmf: "1.0",
  materials: ["PLA"],
  max: [4,4,10.375],
  min: [-4,-4,-0.375],
  units: "mm",
}*/

#define M_PI 3.1415926535897932384626433832795

float coilSquareFace(in mat4 xfm, float radius, float size, float gap, float nTurns, in vec4 xyz) {
  xyz = xyz * xfm;
  // First, constrain the coil to the cylinder with wall thickness "size":
  float rxy = length(xyz.xy);
  if (rxy < (radius-0.5*size) || rxy > (radius + 0.5*size)) { return 0.; }
  // If the current point is between the coils, return no material:
  float angle = atan(xyz.y, xyz.x);
  float z = mod(xyz.z, size+gap) + angle;
  if (z > 0.5*size && z < 0.5*size + gap) { return 0.; }
  // If the current point is within the first coil, stop it at angle < 0.
  float num = floor(xyz.z / (size+gap));
  if (num < 1. && angle < 0.5*M_PI) { return 0.; }
  // Not finished yet...
  return 1.;
}

void mainModel4( out vec4 materials, in vec3 xyz ) {
  materials[0] = coilSquareFace(mat4(1), 4.0, 1.0, 1.0, 10., vec4(xyz,1.));
}
```

* Try loading [coil-1.irmf](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/examples/003-coil-square-face/coil-1.irmf) now in the experimental IRMF editor!

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
