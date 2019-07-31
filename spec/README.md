# IRMF Specification v0.0.1

## Background

An IRMF (“Infinite Resolution Materials Format”) file is a JSON blob containing
(required and optional) key-value pairs followed by a [GLSL ES](https://en.wikipedia.org/wiki/OpenGL_ES#OpenGL_ES_3.0)
shader that is written such that it can “render” (or manufacture) a 3D object at
any resolution desired.
The renderer or 3D printer takes advantage of an on-board GPU to freely
slice the model in any convenient 2D plane and takes as many passes as necessary
to fully define and fabricate the model.

That 2D plane represents the quantity of each material
that the renderer will deposit into 3D space along that 2D plane.
(An unlimited number of homogeneous or full-color materials are supported.)
By modifying the parameters while the printer is building the model and
re-slicing it from different angles and positions, the 3D printer can get
the information it needs to build the model.
Additionally, triplets (or groups) of material values can be combined
to represent a full-color (RGB, HSL, or any format) spectrum for a single material.
There is nothing in the spec that limits the interpretation (or range) of
the material values output by the IRMF shader.

Each material value typically varies from 0 to 1, representing no material up
to solid material. There is no checking that the material values sum up to 1,
which allows the 3D printer manufacturer to use the values in clever and
differentiating ways. But in general, if the material value is less than or
equal to 0, there will be no material, and if the value is greater than
or equal to 1, there will be solid material in that location. (This makes it
much easier to write shaders since you don't need to worry about clamping
the material values explicitly between 0 and 1.)

## Format Specifications

An IRMF file (also known as an IRMF shader) *MUST* start with the following three
characters followed by a `\n` (newline or `\r\n` [carriage return, newline]
on DOS systems):

* `/*{`

Immediately following this opening are JSON key-value pairs
(listed in any order) that describe the properties of the shader.
Here are the keys and sample values:

* `author: "<name of author>",`
  * *optional* - e.g. `"Glenn M. Lewis"`
* `copyright: "<copyright text>",`
  * *optional* - e.g. `"Apache-2.0"`
* `date: "<date created>",`
  * *optional* - e.g. `"2019-06-28"`
* `encoding: "<shader encoding>",`
  * *optional* - this specifies the encoding of the shader portion of the file.
    Possible values are:
    * _null_ or `""` - default (ASCII) encoding.
    * "gpg" - encrypted [GnuPG](https://gnupg.org/) ASCII-armored encoding for a
      particular device or printer manufacturer. This can be useful if someone
      wanted to sell their model but didn't want to release the shader source code.
    * "gzip" - binary compressed [gzip](https://en.wikipedia.org/wiki/Gzip) format.
    * "gzip+base64" - ASCII [Base64](https://en.wikipedia.org/wiki/Base64) encoding
      of [gzip'd](https://en.wikipedia.org/wiki/Gzip) binary compressed format.
      This can be useful when the IRMF shader gets large (for example, when rendering
      text in your model).
* `glslVersion: "#version 300 es",`
  * *optional* - this specifies the GLSL version being used in the shader.
    This defaults to `"#version 300 es"`.
* `irmf: "1.0",`
  * *required* - this is the _major.minor_ version of the IRMF spec.
* `materials: ["<m1 name>","<m2 name>","<m3 name>","<m4 name>"],`
  * *required* - must be the same length as the number of material values
     output by this IRMF shader. e.g. `["support","dielectric","AISI 1018 steel"]`.
     The material name is used to identify the desired material to the 3D printer.
    * For 1-4 materials, the `mainModel4` function will be used.
    * For 5-9 materials, the `mainModel9` function will be used.
    * For 10-16 materials, the `mainModel16` function will be used.
    * For 17-32 materials, the `mainModel32` function will be used.
    * For 33-48 materials, the `mainModel48` function will be used, _etc_.
* `max: [<urx>,<ury>,<urz>],`
  * *required* - upper right bounds of shader - e.g. `[0,0,0]`.
* `min: [<llx>,<lly>,<llz>],`
  * *required* - lower left bounds of shader - e.g. `[10,12,15]`.
* `notes: "<notes from IRMF shader author>",`
  * *optional*.
* `options: {<key1>: <value1>, <key2>: <value2> [,...]},`
  * *optional* - These key-value pairs can be used by the renderer or 3D printer
    as custom options that control the viewing or manufacturing of models.

    They are renderer- (or device-)specific. Renderers or 3D printers that don't
    recognize individual options will simply ignore them (possibly with a warning).

    *e.g.* `{ showAxes: false, showSliders: false, goldPlating: "1um" }`.
* `title: "<name of IRMF model>",`
  * *optional*
* `units: "mm",`
  * *required* - can be `"mm"` or `"in"` or any string that the 3D printer chooses
    to support (*e.g.* `"nm"`), but the units must be specified with the model.
* `version: "<IRMF shader version>",`
  * *optional* - determined by the IRMF shader author - e.g. `"2.7"`.

After the JSON key-value pairs, the following group of three characters *MUST*
be on a line by itself:

* `}*/`

What immediately follows this JSON blob is a GLSL ES shader
([quick reference PDF](https://www.khronos.org/files/opengles3-quick-reference-card.pdf))
(using the specified `encoding`) similar to a ShaderToy
[`mainImage`](https://www.shadertoy.com/howto)
“pixel shader” (or “full-screen fragment shader”), but instead of this
ShaderToy function signature:

* `void mainImage( out vec4 fragColor, in vec2 fragCoord )`

IRMF instead uses one of the following (depending on how many materials
are named in the JSON blob header):

* `void mainModel4( out vec4 materials, in vec3 xyz )`
  (for 1-4 materials)
* `void mainModel9( out mat3 materials, in vec3 xyz )`
  (for 5-9 materials)
* `void mainModel16( out mat4 materials, in vec3 xyz )`
  (for 10-16 materials)
* `void mainModel32( out mat4 materialsA, out mat4 materialsB, in vec3 xyz )`
  (for 17-32 materials)
* `void mainModel48( out mat4 materialsA, out mat4 materialsB, out mat4 materialsC, in vec3 xyz )`
  (for 33-48 materials)

  ... _etc._ ...

The `xyz` input can range anywhere within the minimum bounding box
defined in the JSON blob header. The units are specified in the
header.

The renderer modifies this function on each slice of the design in order
to calculate the amount of material needed at each point in 3D space. It is
free to “zoom in” to any portion of the design to get as much detail as
necessary to generate the model. This is why IRMF shaders have infinite
resolution. The renderer can get as much detail from the shader as it needs
in order to manufacture the part within the alloted timeframe. Higher
resolution models typically take more time to manufacture, so the same
IRMF shader can be used to create quick prototypes or highly-detailed,
final production-worthy parts.
