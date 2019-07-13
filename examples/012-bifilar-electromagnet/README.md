# 012-bifilar-electromagnet

This is the model I wanted to build using convential CAD tools that simply
were not capable of handling the complexity without jumping through hoops
to accomodate the inner workings of the tools (typically by breaking up
the model into much smaller parts that it could more easily handle).

But even after breaking up the model, the CAD tools would attempt to output
STL files to represent the design, which would end up being hundreds of
megabytes (MB). Online 3D printing sites have a maximum upload size limit that
this design exceeded.

So I determined that there must be a better way, and I believe I have
finally found it... IRMF shaders. The next step will be to get IRMF shader
support built into the 3D printers themselves so that these shaders can
be sent directly to the printer as input, and out comes the part as fast
as the printer can make it. No STL. No slicing. No G-Code. Just the
IRMF shader.

## bifilar-electromagnet-1.irmf

```glsl
/*{
  irmf: "1.0",
  materials: ["metal", "dielectric"],
  max: [25,25,61],
  min: [-25,-25,-61],
  units: "mm",
}*/

#define M_PI 3.1415926535897932384626433832795

float coilSquareFace(in mat4 xfm, float radius, float size, float gap, float nTurns, in vec3 xyz) {
  xyz = (vec4(xyz, 1.0) * xfm).xyz;
  
  // First, trivial reject on the two ends of the coil.
  if (xyz.z < -0.5 * size || xyz.z > nTurns * (size + gap) + 0.5 * size) { return 0.0; }
  
  // Then, constrain the coil to the cylinder with wall thickness "size":
  float rxy = length(xyz.xy);
  if (rxy < radius - 0.5 * size || rxy > radius + 0.5 * size) { return 0.0; }
  
  // If the current point is between the coils, return no material:
  float angle = atan(xyz.y, xyz.x) / (2.0 * M_PI);
  if (angle < 0.0) { angle += 1.0; } // 0 <= angle <= 1 between coils
  float dz = mod(xyz.z, size + gap); // 0 <= dz <= (size+gap) between coils.
  
  float lastHelixZ = angle * (size + gap);
  if (lastHelixZ > dz) { lastHelixZ -= (size + gap); }
  float nextHelixZ = lastHelixZ + (size + gap);
  
  if (dz > lastHelixZ + 0.5 * size && dz < nextHelixZ - 0.5 * size) { return 0.0; }
  
  // If the current point is within start of the first coil, stop it at angle < 0.
  if (xyz.z < 0.5 * size && angle > 0.5) { return 0.0; }
  // If the current point is with the end of the last coil, stop it at angle > PI.
  if (xyz.z > nTurns * (size + gap) - 0.5 * size && angle < 0.5) { return 0.0; }
  
  return 1.0;
}

float box(vec3 start, vec3 end, float size, in vec3 xyz) {
  vec3 ll = min(start, end) - vec3(0.5 * size);
  vec3 ur = max(start, end) + vec3(0.5 * size);
  if (any(lessThan(xyz, ll))|| any(greaterThan(xyz, ur))) { return 0.0; }
  return 1.0;
}

 mat3 rotAxis(vec3 axis, float a) {
  // This is from: http://www.neilmendoza.com/glsl-rotation-about-an-arbitrary-axis/
  float s = sin(a);
  float c = cos(a);
  float oc = 1.0 - c;
  vec3 as = axis * s;
  mat3 p = mat3(axis.x * axis, axis.y * axis, axis.z * axis);
  mat3 q = mat3(c, - as.z, as.y, as.z, c, - as.x, - as.y, as.x, c);
  return p * oc + q;
 }

 mat4 rotZ(float degrees) {
  return mat4(rotAxis(vec3(0, 0, 1), M_PI * degrees / 180.0));
 }

float coilPlusConnectorWires(float coilNum, float numCoils, float inc, float innerRadius, float connectorRadius, float size, float gap, float nTurns, in vec3 xyz) {
  mat4 xfm = mat4(1) * rotZ(coilNum * inc);
  float coilRadius = coilNum + innerRadius;
  float coil = coilSquareFace(xfm, coilRadius, size, gap, nTurns, xyz);
  
  vec3 coilXYZ = (vec4(xyz, 1.0) * xfm).xyz;
  
  float bz = -(size + gap);
  float tz = nTurns * (size + gap);
  float tzp1 = (nTurns + 1.0) * (size + gap);
  
  coil += box(vec3(coilRadius, 0.0, 0.0), vec3(coilRadius, 0.0, bz), size, coilXYZ);
  coil += box(vec3(coilRadius, 0.0, bz), vec3(connectorRadius, 0.0, bz), size, coilXYZ);
  coil += box(vec3(connectorRadius, 0.0, bz), vec3(connectorRadius, 0.0, tzp1), size, coilXYZ);
  coil += box(vec3(coilRadius, 0.0, tz), vec3(coilRadius, 0.0, tzp1), size, coilXYZ);
  
  mat4 nextCoilXfm = mat4(1) * rotZ((coilNum + 1.0) * inc);
  vec3 nextCoilXYZ = (vec4(xyz, 1.0) * nextCoilXfm).xyz;
  
  return coil;
}

 vec2 bifilarElectromagnet(float size, float gap, float nTurns, in vec3 xyz) {
  // if (xyz.z < 120.1) { return vec2(0); } // For debugging ends.
  // if (xyz.z > 0.0) { return vec2(0); }  // For debugging ends.
  
  const float numCoils = 20.0;
  const float inc = 360.0 / numCoils;
  const float innerRadius = 3.0;
  float connectorRadius = innerRadius + numCoils * (size + gap);
  
  float coil01 = coilPlusConnectorWires(0.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil02 = coilPlusConnectorWires(1.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil03 = coilPlusConnectorWires(2.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil04 = coilPlusConnectorWires(3.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil05 = coilPlusConnectorWires(4.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil06 = coilPlusConnectorWires(5.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil07 = coilPlusConnectorWires(6.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil08 = coilPlusConnectorWires(7.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil09 = coilPlusConnectorWires(8.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil10 = coilPlusConnectorWires(9.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil11 = coilPlusConnectorWires(10.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil12 = coilPlusConnectorWires(11.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil13 = coilPlusConnectorWires(12.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil14 = coilPlusConnectorWires(13.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil15 = coilPlusConnectorWires(14.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil16 = coilPlusConnectorWires(15.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil17 = coilPlusConnectorWires(16.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil18 = coilPlusConnectorWires(17.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil19 = coilPlusConnectorWires(18.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  float coil20 = coilPlusConnectorWires(19.0, numCoils, inc, innerRadius, connectorRadius, size, gap, nTurns, xyz);
  
  float metal = coil01 + coil02 + coil03 + coil04 + coil05 + coil06 +
  coil07 + coil08 + coil09 + coil10 + coil11 + coil12 + coil13 +
  coil14 + coil15 + coil16 + coil17 + coil18 + coil19 + coil20;
  
  float dielectric = 0.0;
  
  return vec2(metal, dielectric);
 }

 void mainModel4(out vec4 materials, in vec3 xyz) {
  xyz.z += 60.0;
  materials.xy = bifilarElectromagnet(0.85, 0.15, 120.0, xyz);
 }
```

* Try loading [bifilar-electromagnet-1.irmf](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/examples/012-bifilar-electromagnet/bifilar-electromagnet-1.irmf) now in the experimental IRMF editor!

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
