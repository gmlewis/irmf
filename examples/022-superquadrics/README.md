# 022-superquadrics

One of my favorite computer graphics classes I took at Caltech
in 1988 was taught by the teaching assistant,
[John Snyder](https://www.microsoft.com/en-us/research/people/johnsny/)
for the professor, [Al Barr](http://www.gg.caltech.edu/~barr/index.html).

In that class, we wrote our own renderer and one of the primitives
was [superquadrics](https://authors.library.caltech.edu/9756/).

Below are example superquadrics that are beautifully suited for
use in IRMF.

## superquad-ellipsoids-1.irmf

This is a replication of Figure 7 in Al Barr's article linked above.

![superquad-ellipsoids-1.png](superquad-ellipsoids-1.png)

```glsl
/*{
  irmf: "1.0",
  materials: ["PLA"],
  max: [5,5,5],
  min: [-5,-5,-5],
  units: "mm",
}*/

float superquad(in float e1, in float e2, in vec3 xyz) {
  xyz = abs(xyz); // Due to GLSL 'pow' definition.
  float f = pow(pow(xyz.x, 2.0 / e2) + pow(xyz.y, 2.0 / e2), e2 / e1) + pow(xyz.z, 2.0 / e1);
  return f <= 1.0 ? 1.0 : 0.0;
}

void mainModel4(out vec4 materials, in vec3 xyz) {
  float u = 1.2;
  float v = 3.4;
  materials[0] =
  superquad(0.3, 0.3, xyz - vec3(-v, 0, v))
  + superquad(0.3, 0.1, xyz - vec3(-u, 0, v))
  + superquad(0.3, 1.0, xyz - vec3(u, 0, v))
  + superquad(0.3, 3.0, xyz - vec3(v, 0, v))
  + superquad(0.1, 0.3, xyz - vec3(-v, 0, u))
  + superquad(0.1, 0.1, xyz - vec3(-u, 0, u))
  + superquad(0.1, 1.0, xyz - vec3(u, 0, u))
  + superquad(0.1, 3.0, xyz - vec3(v, 0, u))
  + superquad(1.0, 0.3, xyz - vec3(-v, 0, - u))
  + superquad(1.0, 0.1, xyz - vec3(-u, 0, - u))
  + superquad(1.0, 1.0, xyz - vec3(u, 0, - u))
  + superquad(1.0, 3.0, xyz - vec3(v, 0, - u))
  + superquad(3.0, 0.3, xyz - vec3(-v, 0, - v))
  + superquad(3.0, 0.1, xyz - vec3(-u, 0, - v))
  + superquad(3.0, 1.0, xyz - vec3(u, 0, - v))
  + superquad(3.0, 3.0, xyz - vec3(v, 0, - v));
}
```

* Try loading [superquad-ellipsoids-1.irmf](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/examples/022-superquadrics/superquad-ellipsoids-1.irmf) now in the experimental IRMF editor!

* Here is a crude STL approximation of this model
  using [irmf-slicer](https://github.com/gmlewis/irmf-slicer):
  - [superquad-ellipsoids-1-mat01-PLA.stl](superquad-ellipsoids-1-mat01-PLA.stl) (32153684 bytes)

## superquad-ellipsoids-2.irmf

This version is sliced into two materials to make it easier to visualize.

![superquad-ellipsoids-2.png](superquad-ellipsoids-2.png)

```glsl
/*{
  irmf: "1.0",
  materials: ["Red","Green"],
  max: [5,5,5],
  min: [-5,-5,-5],
  units: "mm",
}*/

vec2 superquad2(in float slices, in float e1, in float e2, in vec3 xyz) {
  xyz = abs(xyz); // Due to GLSL 'pow' definition.
  float f = pow(pow(xyz.x, 2.0 / e2) + pow(xyz.y, 2.0 / e2), e2 / e1) + pow(xyz.z, 2.0 / e1);
  if (f > 1.0) { return vec2(0); }
  return (mod(abs(xyz.z) * slices, 1.0) <= 0.5) ? vec2(1, 0) : vec2(0, 1);
}

void mainModel4(out vec4 materials, in vec3 xyz) {
  materials.xy = vec2(0);
  
  vec4 e = vec4(0.1, 1.0, 2.0, 3.0);
  vec4 o = vec4(-3.4, - 1.2, 1.2, 3.4);
  for(int i = 0; i < 4; i ++ ) {
    for(int j = 0; j < 4; j ++ ) {
      materials.xy += superquad2(6.0, e[i], e[j], xyz - vec3(o[j], 0, o[3 - i]));
    }
  }
}
```

* Try loading [superquad-ellipsoids-2.irmf](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/examples/022-superquadrics/superquad-ellipsoids-2.irmf) now in the experimental IRMF editor!

## superquad-toroids-1.irmf

This is a replication of Figure 10 in Al Barr's article linked above.

![superquad-toroids-1.png](superquad-toroids-1.png)

```glsl
/*{
  irmf: "1.0",
  materials: ["PLA"],
  max: [5,5,5],
  min: [-5,-5,-5],
  units: "mm",
}*/

float superquad(in float e1, in float e2, in float a4, in vec3 xyz) {
  xyz.xyz *= 2.5;
  float angle = -1.3;
  float c = cos(angle);
  float s = sin(angle);
  xyz.yz = mat2(c, - s, s, c) * xyz.yz;
  xyz = abs(xyz); // Due to GLSL 'pow' definition.
  float f = pow(pow(pow(xyz.x, 2.0 / e2) + pow(xyz.y, 2.0 / e2), e2 / 2.0) - a4, 2.0 / e1) + pow(xyz.z, 2.0 / e1);
  return f <= 1.0 ? 1.0 : 0.0;
}

void mainModel4(out vec4 materials, in vec3 xyz) {
  float u = 1.2;
  float v = 3.4;
  float a4 = 1.25;
  materials[0] =
  superquad(0.3, 0.3, a4, xyz - vec3(-v, 0, v))
  + superquad(0.3, 0.1, a4, xyz - vec3(-u, 0, v))
  + superquad(0.3, 1.0, a4, xyz - vec3(u, 0, v))
  + superquad(0.3, 3.0, a4, xyz - vec3(v, 0, v))
  + superquad(0.1, 0.3, a4, xyz - vec3(-v, 0, u))
  + superquad(0.1, 0.1, a4, xyz - vec3(-u, 0, u))
  + superquad(0.1, 1.0, a4, xyz - vec3(u, 0, u))
  + superquad(0.1, 3.0, a4, xyz - vec3(v, 0, u))
  + superquad(1.0, 0.3, a4, xyz - vec3(-v, 0, - u))
  + superquad(1.0, 0.1, a4, xyz - vec3(-u, 0, - u))
  + superquad(1.0, 1.0, a4, xyz - vec3(u, 0, - u))
  + superquad(1.0, 3.0, a4, xyz - vec3(v, 0, - u))
  + superquad(3.0, 0.3, a4, xyz - vec3(-v, 0, - v))
  + superquad(3.0, 0.1, a4, xyz - vec3(-u, 0, - v))
  + superquad(3.0, 1.0, a4, xyz - vec3(u, 0, - v))
  + superquad(3.0, 3.0, a4, xyz - vec3(v, 0, - v));
}
```

* Try loading [superquad-toroids-1.irmf](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/examples/022-superquadrics/superquad-toroids-1.irmf) now in the experimental IRMF editor!

* Here is a crude STL approximation of this model
  using [irmf-slicer](https://github.com/gmlewis/irmf-slicer):
  - [superquad-toroids-1-mat01-PLA.stl](superquad-toroids-1-mat01-PLA.stl) (48618084 bytes)

## sphericon-1.irmf

A [sphericon](https://en.wikipedia.org/wiki/Sphericon) is easy to make
out of two half-superquadrics.

![sphericon-1.png](sphericon-1.png)

```glsl
/*{
  irmf: "1.0",
  materials: ["PLA"],
  max: [5,5,5],
  min: [-5,-5,-5],
  units: "mm",
}*/

float superquad(in float e1, in float e2, in vec3 xyz) {
  xyz = abs(xyz); // Due to GLSL 'pow' definition.
  float f = pow(pow(xyz.x, 2.0 / e2) + pow(xyz.y, 2.0 / e2), e2 / e1) + pow(xyz.z, 2.0 / e1);
  return f <= 1.0 ? 1.0 : 0.0;
}

float sphericon(in vec3 xyz) {
  if (xyz.x <= 0.0) {
    return superquad(2.0, 1.0, xyz);
  }
  return superquad(2.0, 1.0, xyz.xzy);
}

void mainModel4(out vec4 materials, in vec3 xyz) {
  xyz /= 5.0;
  materials[0] = sphericon(xyz);
}
```

* Try loading [sphericon-1.irmf](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/examples/022-superquadrics/sphericon-1.irmf) now in the experimental IRMF editor!

* Here is a crude STL approximation of this model
  using [irmf-slicer](https://github.com/gmlewis/irmf-slicer):
  - [sphericon-1-mat01-PLA.stl](sphericon-1-mat01-PLA.stl) (32787084 bytes)

## sphericon-2.irmf

It turns out that it is much easier to visualize a sphericon in two colors
(at least in the IRMF editor), which translates to using two materials.

![sphericon-2.png](sphericon-2.png)

```glsl
/*{
  irmf: "1.0",
  materials: ["Red","Green"],
  max: [5,5,5],
  min: [-5,-5,-5],
  units: "mm",
}*/

float superquad(in float e1, in float e2, in vec3 xyz) {
  xyz = abs(xyz); // Due to GLSL 'pow' definition.
  float f = pow(pow(xyz.x, 2.0 / e2) + pow(xyz.y, 2.0 / e2), e2 / e1) + pow(xyz.z, 2.0 / e1);
  return f <= 1.0 ? 1.0 : 0.0;
}

vec2 sphericon2(in float slices, in vec3 xyz) {
  if (xyz.x <= 0.0) {
    float v = superquad(2.0, 1.0, xyz);
    if (mod(abs(xyz.z) * slices, 1.0) <= 0.5) { return vec2(v, 0); }
    return vec2(0, v);
  }
  float v = superquad(2.0, 1.0, xyz.xzy);
  if (mod(-abs(xyz.y) * slices, 1.0) <= 0.5) { return vec2(v, 0); }
  return vec2(0, v);
}

void mainModel4(out vec4 materials, in vec3 xyz) {
  xyz /= 5.0;
  materials.xy = sphericon2(6.0, xyz);
}
```

* Try loading [sphericon-2.irmf](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/examples/022-superquadrics/sphericon-2.irmf) now in the experimental IRMF editor!

* Here is a crude STL approximation of this model
  using [irmf-slicer](https://github.com/gmlewis/irmf-slicer)
  (one STL file per material):
  - [sphericon-2-mat01-Red.stl](sphericon-2-mat01-Red.stl) (42254484 bytes)
  - [sphericon-2-mat02-Green.stl](sphericon-2-mat02-Green.stl) (42251284 bytes)

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
