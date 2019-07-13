# Infinite Resolution Materials Format (IRMF) Shaders

## Summary

IRMF is a file format used to describe [GLSL
ES](https://en.wikipedia.org/wiki/OpenGL_ES) shaders that define the
materials in a 3D object with infinite resolution. IRMF completely
eliminates the need for [software
slicers](https://en.wikipedia.org/wiki/Slicer_(3D_printing)),
[STL](https://en.wikipedia.org/wiki/STL_(file_format)), and
[G-code](https://en.wikipedia.org/wiki/G-code) files used in
[3D printers](https://en.wikipedia.org/wiki/3D_printing).

I believe that IRMF shaders will revolutionize the 3D-printing industry.

## Introduction

This article is about a revolutionary computational technique,
infinite resolution materials format (IRMF) shaders, that takes
digitally-generated 3D models to the next level. You can think of it
as the equivalent of Gutenberg’s press for the creation of 3D
physical objects.

Historically, geometric primitives such as spheres, cubes, cones,
etc. were combined to build up complex 3D models. Today, sophisticated
parametric CAD programs such as [Onshape](https://cad.onshape.com/)
allow designers to model complex 3D parts. However, it takes a great
deal of computing power to run the algorithms necessary to convert these
models into STL triangle meshes that are then sliced and converted to
G-code that a 3D printer can use to manufacture the objects.
Additionally, extremely complex parts, especially organic
shapes, are not well suited to this style of design and frequently
overwhelm the legacy algorithms and slicing programs. Most importantly,
the resulting STL triangle mesh doesn’t accurately describe the desired
surface of the model, but is only a tessellated approximation of the
truly-desired object, having limited resolution. Not only that, but the
higher the desired resolution, the larger the STL file itself
grows... to the point where it can not be easily handled by the tools
that need to process it.

IRMF shaders turn this operation inside-out. Instead of imperatively
building up an object with a sequence of steps, IRMF shaders take a
functional declarative approach. IRMF shaders answer the question “What
material exists at this point in 3D space?” The advantage to this
approach is that the same shader can instantly provide a voxelized
representation of the object at any desired output resolution by simply
instantiating as many instances of this IRMF shader as needed. This is
what is called an “embarrassingly parallel” computation and is what
[GPUs](https://en.wikipedia.org/wiki/Graphics_processing_unit) were
designed for (although typically for rendering images on a 2D
screen). With IRMF shaders, GPUs are used to determine what material is
printed at each point in 3D space.

## What is an “infinite resolution materials format (IRMF) shader?”

An IRMF shader consists of two parts: a
[JSON](https://en.wikipedia.org/wiki/JSON) blob description and a set of
instructions used to determine what material is placed at any point in
3D space. What makes it a shader is that it is used simultaneously for
every voxel (volumetric pixel) within the object. This means that the
code has to behave differently for every (x,y,z) position within 3D
space. Like a type press, the shader is passed a position in space and
returns the percentages of one or more materials. When compiled and run
in parallel (one shader per point in 3D space), it will be incredibly
fast... nearly instantaneous.

If the GPU does not have enough capacity to assign one shader per
point in 3D space, the design can be diced up into cubes or slices
and then reassembled, or more GPUs could be employed to cover the
full 3D object at the resolution desired.

Alternatively, depending on the style of 3D printer, the IRMF
shader could be processed within the 3D printer itself and no slicing
step would be needed at all! Light-based 3D printers are
exceptionally well suited to this paradigm. As the 3D printer is
ready to create material, it simply asks the IRMF shader what
percentage of materials belong at each location, and generates the
materials without the need for slicing at all. The 3D printer could
simply accept the extremely compact IRMF file itself as input, then
generate the 3D object at any resolution the printer supports. In
fact, one option on the 3D printer might be “How fast do you want
this part? It can be made with 1μm resolution in 1 minute and 1nm
resolution in 10 minutes.” The exact same IRMF file describes the
model with infinite resolution and does not need to be changed in
this scenario.

The JSON blob is used to describe the physical dimensions (the
minimum bounding box) of the object, how many materials it uses, and
other parameters used by the shader. The shader portion itself is a
standard GLSL ES shader.

## Do IRMF shaders really have infinite resolution?

Well, yes and no. They are currently limited by the resolution of a
GPU's `float` representation. However, the IRMF shader itself is pure
math and does not limit the resolution of the model, so `IRMF` is truly
an appropriate term for the shaders as they are not the limiting factor.

When GPUs have higher resolution in their numeric representations,
no major changes will be needed to update the IRMF shaders... most
likely it will involve a simple name change from `float` to the new
keyword.

## Inspiration

[Shadertoy.com](https://shadertoy.com) is an amazing collection
of GLSL ES shaders written by a lot of amazing, creative people. From
there, I found this incredible website: [The Book of
Shaders](https://thebookofshaders.com/) which teaches shader writing
from the ground up.

Additionally, I came across a similar use of JSON and GLSL ES called
[ISF](https://www.interactiveshaderformat.com/) but used for video.

## What is the difference between GLSL ES and IRMF?

GLSL ES is the [OpenGL Shading
Language](https://www.khronos.org/opengles/) developed by the [Khronos
Group](https://www.khronos.org/). GLSL ES files are compiled into
shaders that can be run in parallel on GPU cards.

IRMF is designed to be a standard for working with GLSL ES in such a
way that 3D printers can manufacture objects at any resolution
possible. IRMF files consist of a JSON blob followed by a GLSL ES
shader.

## How are IRMF shaders different from signed distance functions (SDFs)?

IRMF shaders are *much* easier to write than [SDFs](https://github.com/gmlewis/sdfx).
In fact, one could write [genetic programs](https://github.com/gmlewis/gep)
to generate IRMF shaders.

Back in 2018, I was fed up with STL models and CAD tools that couldn't handle
non-trivial booleans (think [bifilar coils](https://github.com/gmlewis/go-gerber))
so I came across voxels and [wrote some tools](https://github.com/gmlewis/stldice)
to make it easier to manipulate voxel designs and perform complex boolean operations
that cause all other popular CAD tools to choke.

Then I came to the realization that voxels, although they solved the boolean operation
problems, did not address the problem of generating smooth curves and surfaces
because they are inherently limited by the resolution of the image slices.

Finally, I came across [SDFs](https://github.com/gmlewis/sdfx) and thought that I
had found the modeling tool that would solve all the boolean and smooth surface
problems for good. But then I contributed my first primitive, a [spiral](
https://github.com/gmlewis/sdfx/blob/master/sdf/spiral.go), and found it incredibly
difficult to get it right.

In a 2D or 3D SDF, you must specify the signed *distance* from any point in 2D or 3D
space (respectively) to the *nearest* edge of your primitive (negative is inside the
object, zero is the edge, and positive is outside the object).
Think of implementing this for a spiral.
That is a royal pain in the rear. The math was not fun.

An IRMF shader, on the other hand, is given a single point in 3D space and all that
is needed is for the shader to report what material(s) exist in that one point in
space, and not how far it is to some other material. For the spiral case, you only
need to determine if the point is inside the spiral or outside it, not how far the
point is to the edge.

This makes IRMF shaders *orders of magnitude* easier to write than SDFs.
IRMF shaders make boolean operations a breeze: zero times anything is zero; one
times anything is that thing. Booleans solved. Likewise, IRMF shaders can represent
curved surfaces as fine as the GPU can resolve... which is mighty fine.

I'm very excited about the future of IRMF shaders, and envision a day when the
Star Trek replicator will be a household device. “Hey replicator, print me a widget.”
“OK, your widget is ready.”

## How do I use IRMF?

An [IRMF shader editor](https://github.com/gmlewis/irmf-editor) is in the works.
Eventually, firmware for 3D printers will be written that natively read, parse,
and use IRMF files to generate physical objects with one or more materials...
thereby completely eliminating the need for STL files, software slicers,
and G-Code files.

## Who created IRMF?

The IRMF Specification was created and is maintained by Glenn M. Lewis.

Issues can be raised on the [GitHub issues page](https://github.com/gmlewis/irmf/issues)
for the IRMF Specification.

## What is the status of IRMF?

Currently, IRMF is just an idea that needs fully fleshing out.

Please see the [IRMF Spec](spec) and [provided examples](examples) for more information.

## Examples

* [001-sphere](examples/001-sphere)
* [002-cube](examples/002-cube)
* [003-coil-square-face](examples/003-coil-square-face)
* [004-coil-circle-face](examples/004-coil-circle-face)
* [005-cylinder](examples/005-cylinder)
* [006-square-tetrahedron](examples/006-square-tetrahedron)
* [007-cone](examples/007-cone)

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
