# 006-square-tetrahedron

## tetrahedron-1.irmf

If we make a tetrahedron with a square base that fits with a (-1,-1,-1)-(1,1,1)
cube, then it could be scaled to any size independently along each axis.

![tetrahedron-1.png](tetrahedron-1.png)

```glsl
/*{
  irmf: "1.0",
  materials: ["PLA"],
  max: [1,1,0.5],
  min: [-1,-1,-0.5],
  units: "mm",
}*/

float squareTetrahedron(in vec3 xyz) {
  // Trivially reject above and below the tetrahedron.
  if (xyz.z < 0.0 || xyz.z > 1.0) { return 0.0; }
  
  // Calculate the new size based on the height.
  float zsize = mix(1.0, 0.0, xyz.z);
  
  if (any(greaterThan(xyz.xy, vec2(zsize)))|| any(lessThan(xyz.xy, vec2(-zsize)))) { return 0.0; }
  
  return 1.0;
}

void mainModel4(out vec4 materials, in vec3 xyz) {
  xyz.z += 0.5;
  materials[0] = squareTetrahedron(xyz);
}
```

* Try loading [tetrahedron-1.irmf](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/examples/006-square-tetrahedron/tetrahedron-1.irmf) now in the experimental IRMF editor!

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
