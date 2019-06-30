# 002-cube

While the sphere is one of the easiest IRMF shaders to write, the cube is actually simpler.

For a cube, we can exploit the fact that the shader values are only valid within the
confines of the minimum bounding box (MBB). Since the MBB of a cube is the cube itself,
we simply need to return a material value of 1 for all values passed to the shader,
and the MBB itself defines the object (the cube).

Here is an [IRMF shader](cube.irmf) defining a 10mm diameter cube:

/*{
  author: "Glenn M. Lewis",
  copyright: "Apache-2.0",
  date: "2019-06-30",
  irmf: "1.0",
  materials: ["PLA"],
  max: [10,10,10],
  min: [0,0,0],
  notes: "Simple IRMF shader - Hello, Cube!",
  title: "10mm diameter Cube",
  units: "mm",
  version: "1.0"
}*/

void mainModel4( out vec4 materials, in vec3 xyz ) {
  materials[0] = 1.0;
}

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
