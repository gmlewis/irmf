# IRMF Specification v0.0.1

## Background

An IRMF ("Infinite Resolution Materials Format") file is a JSON blob containing
required and optional fields followed by a GLSL ES shader that is written such
that it can render a 2D slice of a 3D object at any resolution desired on a GPU.

That 2D slice represents the quantity of up to four materials that the 3D printer
will deposit into 3D space along that 2D plane. By modifying the parameters while
the printer is building the model and re-slicing the model from different angles
and positions, the 3D printer can generate all the information it needs to build
the model using up to four materials. (Future versions of this spec may support
more than four materials.) Additionally, three material values can be combined to
represent a full color spectrum for a single material.

Each material value (located in the R, G, B, and A channels of the fragment shader)
varies from 0 to 1, representing no material up to solid material. There is no
checking that the material values sum up to 1, which allows the 3D printer
manufacturer to use the values in clever ways.

## Format Specifications

An IRMF file (also known as an IRMF shader) *MUST* start with the following three
characters followed by a newline or (carriage return, newline on DOS systems):

* `/*{`

Immediately following this opening are JSON key-value pairs
(listed in any order) that describe the properties of the shader.
Here are the keys and sample values:

* `author: "<name of author>",`
  * (*optional* - e.g. `"Glenn M. Lewis"`)
* `copyright: "<copyright text>",`
  * (*optional* - e.g. `"Apache-2.0"`)
* `date: "<date created>",`
  * (*optional* - e.g. `"2019-06-28"`)
* `irmf: "1.0",`
  * (*required* - this is the version of the IRMF spec)
* `materials: ["<m1 name>","<m2 name>","<m3 name>","<m4 name>"],`
  * (*required* - must be the same length as the number of material values
     output by this IRMF shader. e.g. `["support","AISI 1018 steel"]`)
* `max: [<urx>,<ury>,<urz>],`
  * (*required* - upper right bounds of shader - e.g. `[0,0,0]`)
* `min: [<llx>,<lly>,<llz>],`
  * (*required* - lower left bounds of shader - e.g. `[10,12,15]`)
* `notes: "<notes from IRMF shader author>",`
  * (*optional*)
* `units: "mm",`
  * (*required* - can be `"mm"` or `"in"`)
* `version: "<IRMF shader version>",`
  * (*optional* - e.g. `"2.7"`)

After the JSON key-value pairs, the following group of three characters *MUST*
be on a line by itself:

* `}*/`

What follows this JSON blob is an almost-standard ShaderToy (GLSL ES)
"pixel shader" (or "full-screen fragment shader") with the exception that
a function is provided by the IRMF viewer software or 3D printer (the
*"renderer"*) that transforms the input `in vec2 fragCoord` to a `vec3 xyz`
that (evntually) fully covers the minimum bounding box of the design in the
provided units (typically "mm").

The renderer modifies this function on each slice of the design in order
to calculate the amount of material needed at each point in space. It is
free to "zoom in" to any portion of the design to get as much detail as
necessary to generate the model. This is why IRMF shaders have infinite
resolution. The renderer can get as much detail from the shader as it needs
in order to manufacture the part within the alloted timeframe. Higher
resolution models typically take more time to manufacture, so the same
IRMF shader can be used to create quick prototypes or highly detailed
final production-worthy parts.
