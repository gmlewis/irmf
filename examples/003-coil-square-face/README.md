# 003-coil-square-face

Another surprisingly-simple model is a helical coil with a square cross-section
face.

## coil-1.irmf

```glsl
/*{
  irmf: "1.0",
  materials: ["PLA"],
  max: [4,4,120.375],
  min: [-4,-4,-0.375],
  notes: "Simple IRMF shader - coil with square cross-section face.",
  title: "8mm diameter Coil",
  units: "mm"
}*/

float coilSquareFace(in mat4 xfm, float radius, float size, float gap, float nTurns, in vec4 xyz) {
  // TODO
}

void mainModel4( out vec4 materials, in vec3 xyz ) {
  materials[0] = coilSquareFace(mat4(), 4.0, 0.75, 0.25, 120., vec4(xyz,1.));
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
