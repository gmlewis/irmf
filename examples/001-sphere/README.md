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
  materials: ["dielectric","AISI 1018 steel"],
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
