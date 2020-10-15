# 026-utron

Otis T. Carr invented the utron. Here is a simple model of half of a utron.

## half-utron-1.irmf

```glsl
/*{
  irmf: "1.0",
  materials: ["copper"],
  max: [36,36,46],
  min: [-36,-36,0],
  units: "mm",
}*/

float cone(in vec3 xyz) {
  // Trivially reject above and below the cone.
  if (xyz.z < 0.0 || xyz.z > 1.0) { return 0.0; }
  
  // Calculate the new size based on the height.
  float zsize = mix(1.0, 0.0, xyz.z);
  float r = length(xyz.xy);
  
  if (r > zsize) { return 0.0; }
  
  return 1.0;
}

float sphere(in float radius, in vec3 xyz) {
  float r = length(xyz);
  return r <= radius ? 1.0 : 0.0;
}

float cylinder(float radius, float height, in vec3 xyz) {
  // First, trivial reject on the two ends of the cylinder.
  if (xyz.z < 0.0 || xyz.z > height) { return 0.0; }
  
  // Then, constrain radius of the cylinder:
  float rxy = length(xyz.xy);
  if (rxy > radius) { return 0.0; }
  
  return 1.0;
}

float halfUtron(in float sphereDiam, in float edgeLen, in float shaftDiam, in float shaftHeight, in vec3 xyz) {
    float h = sqrt(edgeLen*edgeLen/2.0);
    float r = sphereDiam/2.0;
    float shaftR = shaftDiam/2.0;
    vec3 shaftPos = vec3(0,0,h-shaftR);
    return cone(xyz/h) + cylinder(shaftR, shaftHeight, xyz-shaftPos) - sphere(1.0, xyz/r);
}

void mainModel4(out vec4 materials, in vec3 xyz) {
  materials[0] = halfUtron(44.0, 50.0, 4.98, 10.0, xyz);
}
```

* Try loading [half-utron-1.irmf](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/examples/026-utron/half-utron-1.irmf) now in the experimental IRMF editor!

----------------------------------------------------------------------

# License

Copyright 2020 Glenn M. Lewis. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
