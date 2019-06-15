# Infinite Resolution Materials Format (IRMF) Shaders
  An invention by Glenn M. Lewis - 2019-03-10

## Summary

IRMF is a file format used to describe GLSL shaders that define
the materials in a 3D object with infinite resolution. IRMF
completely eliminates the need for software slicers, STL, and G-code
files used in 3D printers.

## Introduction

This article is about a revolutionary computational technique,
infinite resolution materials format (IRMF) shaders, that takes
digitally-generated 3D models to the next level. You can think of it
as the equivalent of Gutenberg's press for the creation of 3D
objects.

Historically, geometric primitives such as spheres, cubes, cones,
etc. were combined to build up complex 3D models. Today,
sophisticated parametric CAD programs such as Onshape.com allow
designers to model complex 3D parts. However, it takes a great deal
of computing power to run the algorithms necessary to convert these
models into STL triangle meshes that are then sliced and converted to
G-code that a 3D printer can use to manufacture the
objects. Additionally, extremely complex parts, especially organic
shapes, are not well suited to this style of design and frequently
overwhelm the legacy algorithms and slicing programs. Most
importantly, the resulting STL triangle mesh doesn't accurately
describe the desired surface of the model, but is only a tessellated
approximation of the truly-desired object, having limited
resolution. Not only that, but the higher the desired resolution, the
larger the STL file itself grows... to the point where it can not be
easily handled by the tools that need to process it.

IRMF shaders turn this operation inside-out. Instead of
imperatively building up an object with a sequence of steps, IRMF
shaders take a functional declarative approach. IRMF shaders answer
the question "What material exists at this point in 3D space?" The
advantage to this approach is that the same shader can instantly
provide a voxelized representation of the object at any desired
output resolution by simply instantiating as many instances of this
IRMF shader as needed. This is what is called an "embarrassingly
parallel" computation and is what GPUs were designed for (although
typically for rendering images on a 2D screen). With IRMF shaders,
GPUs are used to determine what material is printed at each point in
3D space.

## What is an "infinite resolution materials format (IRMF) shader"?

An IRMF shader consists of two parts: a JSON blob description and
a set of instructions used to determine what material is placed at
any point in 3D space. What makes it a shader is that it is used
simultaneously for every voxel (volumetric pixel) within the
object. This means that the code has to behave differently for every
(x,y,z) position within 3D space. Like a type press, the shader is
passed a position in space and returns the percentages of one or more
materials. When compiled and run in parallel (one shader per point in
3D space), it will be incredibly fast... nearly instantaneous.

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
fact, one option on the 3D printer might be "How fast do you want
this part? It can be made with 1μm resolution in 1 minute and 1nm
resolution in 10 minutes." The exact same IRMF file describes the
model with infinite resolution and does not need to be changed in
this scenario.

The JSON blob is used to describe the physical dimensions (the
minimum bounding box) of the object, how many materials it uses, and
other parameters used by the shader. The shader portion itself is a
standard GLSL shader.

## Inspiration

[Shadertoy.com](https://shadertoy.com) is an amazing collection
of GLSL shaders written by a lot of amazing creative people. From
there, I found this incredible website: [The Book of
Shaders](https://thebookofshaders.com/) which teaches shader writing
from the ground up.

Additionally, I came across a similar use of JSON and GLSL called
ISF located here: https://www.interactiveshaderformat.com/.

## What is the difference between GLSL and IRMF?

GLSL is the OpenGL Shading Language developed by the Khronos
Group. GLSL files are compiled into shaders that can be run in
parallel on GPU cards.

IRMF is designed to be a standard for working with GLSL in such a
way that 3D printers can manufacture objects at any resolution
possible. IRMF files consist of a JSON blob followed by a GLSL
shader.

## How do I use IRMF?

An IRMF editor is in the works. Eventually, firmware for 3D
printers will be written that natively read, parse, and use IRMF
files to generate physical objects with one or more materials.

## Who created IRMF?

The IRMF Specification was created and is maintained by
Glenn M. Lewis. Issues can be raised on the GitHub page for the
IRMF Specification.

## What is the status of IRMF?

Currently, IRMF is just an idea that needs fully fleshing out.

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
