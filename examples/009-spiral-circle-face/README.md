# 009-spiral-circle-face

## spiral-circle.irmf

Once the square-face spiral was done, it is a rather simple matter
to round the edges.

![spiral-circle.png](spiral-circle.png)

```glsl
/*{
  irmf: "1.0",
  materials: ["PLA"],
  max: [5.5,5.5,0.5],
  min: [-5.5,-5.5,-0.5],
  units: "mm",
}*/

#define M_PI 3.1415926535897932384626433832795

float spiralCircleFace(float startRadius, float size, float gap, float nTurns, in vec3 xyz) {
  // First, trivial reject above and below the spiral.
  if (xyz.z < -0.5 * size || xyz.z > 0.5 * size) { return 0.0; }
  
  float r = length(xyz.xy);
  if (r < startRadius - 0.5 * size ||
      r > startRadius + 0.5 * size + (size + gap) * nTurns) { return 0.0; }
  
  // If the current point is between the spirals, return no material:
  float angle = atan(xyz.y, xyz.x) / (2.0 * M_PI);
  if (angle < 0.0) { angle += 1.0; } // 0 <= angle <= 1 between spirals from center to center.
  // 0 <= dr <= (size+gap) between spirals from center to center.
  float dr = mod(r - startRadius, size + gap);
  
  float coilNum = 0.0;
  float lastSpiralR = angle * (size + gap);
  if (lastSpiralR > dr) {
    lastSpiralR -= (size + gap);  // center of current coil.
    coilNum = -1.0;
  }
  float nextSpiralR = lastSpiralR + (size + gap);  // center of next outer coil.
  
  // If the current point is within the gap between the two coils, reject it.
  if (dr > lastSpiralR + 0.5 * size && dr < nextSpiralR - 0.5 * size) { return 0.0; }
  
  coilNum += floor((r - startRadius + (0.5 * size) - lastSpiralR) / (size + gap));

  // If the current point is in a coil numbered outside the current range, reject it.
  if (coilNum < 0.0 || coilNum >= nTurns) { return 0.0; }

    // At this point, we are within the circle cross-section face, so let's round the edge.
  lastSpiralR = floor((r - startRadius) / (size + gap)) * (size + gap) +
    startRadius + angle * (size + gap);
  if (lastSpiralR > r) { lastSpiralR -= (size + gap); }
  nextSpiralR = lastSpiralR + size + gap;
  angle *= 2.0 * M_PI;
  vec3 lastSpiralCenter = vec3(lastSpiralR * cos(angle), lastSpiralR * sin(angle), 0.0);
  vec3 nextSpiralCenter = vec3(nextSpiralR * cos(angle), nextSpiralR * sin(angle), 0.0);
  float d1 = length(xyz - lastSpiralCenter);
  float d2 = length(xyz - nextSpiralCenter);
  if (d1 > 0.5 * size && d2 > 0.5 * size) { return 0.0; }
  
  return 1.0;
}

void mainModel4(out vec4 materials, in vec3 xyz) {
  materials[0] = spiralCircleFace(3.0, 0.85, 0.15, 2.0, xyz);
}
```

* Try loading [spiral-circle.irmf](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/examples/009-spiral-circle-face/spiral-circle.irmf) now in the experimental IRMF editor!

* Here is a crude STL approximation of this model
  using [irmf-slicer](https://github.com/gmlewis/irmf-slicer):
  - [spiral-circle-mat01-PLA.stl](spiral-circle-mat01-PLA.stl) (47967684 bytes)

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
