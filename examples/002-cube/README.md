# 002-cube

While the sphere is one of the easiest IRMF shaders to write, the cube is actually simpler.

For a cube, we can exploit the fact that the shader values are only valid within the
confines of the minimum bounding box (MBB). Since the MBB of a cube is the cube itself,
we simply need to return a material value of 1 for all values passed to the shader,
and the MBB itself defines the object (the cube).

## cube-1.irmf

Here is an [IRMF shader](cube-1.irmf) defining a 10mm diameter cube:

```glsl
/*{
  irmf: "1.0",
  materials: ["PLA"],
  max: [5,5,5],
  min: [-5,-5,-5],
  units: "mm",
}*/

void mainModel4(out vec4 materials, in vec3 xyz) {
  materials[0] = 1.0;
}
```

* Try loading [cube-1.irmf](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/examples/002-cube/cube-1.irmf) now in the experimental IRMF editor!

## cube-2.irmf

You would probably never write a shader like `cube-1.irmf`. It is just showing
how the minimum bounding box of the shader defines the extent of the model.

To be useful, we would want a `cube` function that could be easily positioned,
rotated, and sized.

Whenever something can be positioned, rotated, and sized, a common
way to do so is to provide a `mat4` matrix that defines all these
transformations in a single bundle. However, to start off, let's explicitly
provide a `pos`ition and `size`.

One thing that the "Book of Shaders" stresses is that the coordinate system
is transformed such that the shader always performs its calculations in its
own local coordinate system. Here's an example of this:

```glsl
/*{
  irmf: "1.0",
  materials: ["PLA"],
  max: [5,5,5],
  min: [-5,-5,-5],
  units: "mm",
}*/

float cube(in vec3 pos, in float size, in vec3 xyz) {
  xyz -= pos; // Move local coordinate system.
  xyz /= size; // Scale local coordinate system.
  if (any(greaterThan(abs(xyz), vec3(0.5)))) {return 0.0; }
  return 1.0;
}

void mainModel4(out vec4 materials, in vec3 xyz) {
  materials[0] = cube(vec3(0), 10.0, xyz);
}
```

* Try loading [cube-2.irmf](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/examples/002-cube/cube-2.irmf) now in the experimental IRMF editor!

## cube-3.irmf

`cube-2.irmf` has the drawback that it can't be rotated. Let's make a
more general-purpose cube-like object that can be translated, scaled,
and rotated.

```glsl
/*{
  irmf: "1.0",
  materials: ["PLA"],
  max: [0.5,0.5,0.5],
  min: [-0.5,-0.5,-0.5],
  units: "mm",
}*/

float cube(in mat4 xfm, in vec4 xyz) {
  xyz = xyz * xfm;
  if (any(greaterThan(abs(xyz), vec4(0.5, 0.5, 0.5, 1.0)))) {return 0.0; }
  return 1.0;
}

void mainModel4(out vec4 materials, in vec3 xyz) {
  materials[0] = cube(mat4(1), vec4(xyz, 1.0)); // mat4(1) is the identity matrix.
}
```

* Try loading [cube-3.irmf](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/examples/002-cube/cube-3.irmf) now in the experimental IRMF editor!

## cube-csg.irmf

```glsl
/*{
  irmf: "1.0",
  materials: ["AISI 1018 steel"],
  max: [5,5,5],
  min: [-5,-5,-5],
  units: "mm",
}*/

float sphere(in vec3 pos, in float radius, in vec3 xyz) {
  xyz -= pos; // Move sphere into place.
  float r = length(xyz);
  return r <= radius ? 1.0 : 0.0;
}

void mainModel4(out vec4 materials, in vec3 xyz) {
  const float radius = 6.0; // 12mm diameter sphere.
  materials[0] = 1.0 - sphere(vec3(0), radius, xyz); // 1.0 is a cube.
}
```

* Try loading [cube-csg.irmf](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/examples/002-cube/cube-csg.irmf) now in the experimental IRMF editor!

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
