# IRMF Specification v1.0

## Background

An IRMF (“Infinite Resolution Materials Format”) file is a JSON blob containing
(required and optional) key-value pairs followed by a [GLSL ES](https://en.wikipedia.org/wiki/OpenGL_ES#OpenGL_ES_3.0),
[WGSL](https://www.w3.org/TR/WGSL/), or [TSL](https://threejs.org/docs/TSL.html)
shader that is written such that it can “render” (or manufacture) a 3D object at
any resolution desired. (GLSL ES and WGSL bodies are shader text; a TSL body is
JavaScript that builds a three.js node graph, which the host lowers to WGSL or
GLSL.) The renderer or 3D printer takes advantage of an on-board GPU to freely
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

* `"author": "<name of author>",`
  * *optional* - e.g. `"Glenn M. Lewis"`
* `"license": "<license text>",`
  * *optional* - e.g. `"Apache-2.0"`
* `"date": "<date created>",`
  * *optional* - e.g. `"2019-06-28"`
* `"encoding": "<shader encoding>",`
  * *optional* - this specifies the encoding of the shader portion of the file.
    Possible values are:
    * _null_ or `""` - default (UTF-8) encoding.
    * "gpg" - encrypted [GnuPG](https://gnupg.org/) ASCII-armored encoding for a
      particular device or printer manufacturer. This can be useful if someone
      wants to sell a model but doesn't want to release the shader source code.
    * "gzip" - binary compressed [gzip](https://en.wikipedia.org/wiki/Gzip) format.
    * "gzip+base64" - ASCII [Base64](https://en.wikipedia.org/wiki/Base64) encoding
      of [gzip'd](https://en.wikipedia.org/wiki/Gzip) binary compressed format.
      This can be useful when the IRMF shader gets large (for example, when rendering
      text in your model).
* `"glslVersion": "#version 300 es",`
  * *optional* - this specifies the GLSL version being used in the shader.
    This defaults to `"#version 300 es"`. This key applies only when
    `"language"` is `"glsl"`; consumers *MUST* ignore it for other languages.
* `"language": "<shader language>",`
  * *required* - this specifies the shader language used. Values are
    case-sensitive. Possible values are:
    * "glsl" - GLSL ES is the shader language; the body is shader text.
    * "wgsl" - WGSL is the shader language; the body is shader text.
    * "tsl" - TSL (Three.js Shading Language) is the shader language; the body
      is a JavaScript program that builds a TSL node graph.
    * A consumer that encounters an unknown `language` value *MUST* reject the
      file with a clear “unsupported language” error. It *MUST NOT* silently
      interpret the body as another language.
* `"irmf": "1.0",`
  * *required* - this is the _major.minor_ version of the IRMF spec.
* `"materials": ["<m1 name>","<m2 name>","<m3 name>","<m4 name>"],`
  * *required* - must be the same length as the number of material values
     output by this IRMF shader. e.g. `["support","dielectric","AISI 1018 steel"]`.
     The material name is used to identify the desired material to the 3D printer.
    * For 1-4 materials, the `mainModel4` function will be used.
    * For 5-9 materials, the `mainModel9` function will be used.
    * For 10-16 materials, the `mainModel16` function will be used.
    * For 17-32 materials, the `mainModel32` function will be used.
    * For 33-48 materials, the `mainModel48` function will be used, _etc_.
* `"max": [<urx>,<ury>,<urz>],`
  * *required* - upper right bounds of shader - e.g. `[0,0,0]`.
* `"min": [<llx>,<lly>,<llz>],`
  * *required* - lower left bounds of shader - e.g. `[10,12,15]`.
* `"notes": "<notes from IRMF shader author>",`
  * *optional*.
* `"options": {<key1>: <value1>, <key2>: <value2> [,...]},`
  * *optional* - These key-value pairs can be used by the renderer or 3D printer
    as custom options that control the viewing or manufacturing of models.

    They are renderer- (or device-)specific. Renderers or 3D printers that don't
    recognize individual options will simply ignore them (possibly with a warning).

    *e.g.* `{ showAxes: false, showSliders: false, goldPlating: "1um" }`.
* `"title": "<name of IRMF model>",`
  * *optional*
* `"tslVersion": "0.180.0",`
  * *optional* - for `"language": "tsl"`, the exact three.js release the shader
    targets (e.g. `"0.180.0"`). It is advisory: consumers *SHOULD* warn (but not
    fail) when it differs from their pinned three.js version, and *SHOULD*
    document the version they pin.
* `"units": "mm",`
  * *required* - can be `"mm"` or `"in"` or any string that the 3D printer chooses
    to support (*e.g.* `"nm"`), but the units must be specified with the model.
* `"version": "<IRMF shader version>",`
  * *optional* - determined by the IRMF shader author - e.g. `"2.7"`.

After the JSON key-value pairs, the following group of three characters *MUST*
be on a line by itself:

* `}*/`

What immediately follows this JSON blob is one of:

* a GLSL ES
  ([GLSL ES 3 quick reference PDF](https://www.khronos.org/files/opengles3-quick-reference-card.pdf))
  shader,
* a WGSL ([WGSL quick reference](https://webgpu.rocks/)) shader, or
* a TSL JavaScript program,

(using the specified `encoding`). The shader text (GLSL/WGSL) is similar to a
ShaderToy [`mainImage`](https://www.shadertoy.com/howto) “pixel shader” (or
“full-screen fragment shader”), but instead of this ShaderToy function
signature:

* `void mainImage( out vec4 fragColor, in vec2 fragCoord )`

IRMF instead defines one of the entry points below, chosen by the number of
materials named in the JSON blob header.

The `xyz` input can range anywhere within the minimum bounding box
defined in the JSON blob header. The units are specified in the
header.

### GLSL ES entry points

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

### WGSL entry points

* `fn mainModel4(xyz: vec3f) -> vec4f` (for 1-4 materials)
* `fn mainModel9(xyz: vec3f) -> mat3x3f` (for 5-9 materials)
* `fn mainModel16(xyz: vec3f) -> mat4x4f` (for 10-16 materials)

Components are packed in material order. For `vec4f` the order is
`.x`, `.y`, `.z`, `.w`. For `mat3x3f`/`mat4x4f` the order is column-major
(i.e. `m[0][0]`, `m[0][1]`, `m[0][2]`, `m[1][0]`, ...), matching the GLSL
`out mat3`/`out mat4` convention. Unused components are ignored.

### TSL entry points

A TSL body is a JavaScript program. The host evaluates it once with the names
exported by `three/tsl` available as lexical bindings (`Fn`, `float`, `vec2`,
`vec3`, `vec4`, `mat2`, `mat3`, `mat4`, `If`, `Loop`, `select`, `uniform`,
`positionWorld`, `instanceIndex`, `texture`, `array`, `Var`, ...). The
`three/tsl` namespace itself *MAY* also be available as `TSL`.

Depending on the number of materials, the body *MUST* define one of:

* `mainModel4` (for 1-4 materials)
* `mainModel9` (for 5-9 materials)
* `mainModel16` (for 10-16 materials)

Each entry point is a top-level callable that accepts one argument — a `vec3`
node `xyz` in model coordinates within the header's `min`…`max` box — and
returns an array of exactly `materials.length` scalar float nodes, in material
order. The idiomatic form is `Fn`:

```js
const mainModel4 = Fn(([xyz]) => {
  // build nodes from xyz
  return [m0, m1]; // one node per material
});
```

A plain function works too: `const mainModel4 = (xyz) => [ /* nodes */ ];`.
Hosts *MAY* additionally accept a `vec4`/`mat3`/`mat4` node for
1–4/5–9/10–16 materials and unpack it in the same order.

The returned values have the same meaning as the GLSL/WGSL outputs: each value
is the amount of that material at `xyz` (typically 0…1), and triplets named
`*.H`/`*.S`/`*.V`, `*.H`/`*.S`/`*.L`, or `*.R`/`*.G`/`*.B` represent full-color
channels exactly as in GLSL/WGSL. The shader *MUST NOT* return an RGB `Color`
node in place of material values; the host performs color mixing.

### Language-specific requirements

**Execution model.** `xyz` is a node, not a concrete coordinate. The graph is
built once and evaluated by the GPU for every sample. Per-sample branching
*MUST* use `select`/`If`; loops *MUST* use `Loop`/`While`. A host is not
required to call the entry point once per voxel. (For GLSL/WGSL the equivalent
is: the shader runs once per sample and computes all material values from
`xyz`.)

**TSL purity and determinism.** A TSL body *MUST NOT* contain `import` or
`export` statements and *MUST NOT* use top-level `await` (hosts *MAY* accept
and strip a leading `import { ... } from 'three/tsl'`). It *MUST NOT* access the
DOM, network, filesystem, `Date`, `Math.random`, or other nondeterministic or
side-effecting APIs. Two evaluations of the same body with the same inputs
*MUST* produce the same graph and the same material values.

**TSL module resolution.** Canonical TSL bodies are self-contained and rely on
the host-injected `three/tsl` names rather than `import` statements, so that
the same body runs in a browser and in a Node-based compiler without import
maps.

**Security.** TSL bodies are executable JavaScript. Consumers *MUST* treat
every `.irmf` file as untrusted code until the user has decided otherwise. A
consumer that evaluates TSL *MUST* run it in a sandbox that provides no network
access, no filesystem access outside the consumer's own read-only resources, no
access to the host environment (including environment variables and
credentials), and enforced time and memory limits.

In an interactive host, remote TSL *MUST* require explicit per-source consent
that names the origin (and, ideally, shows a content hash). Consent *MUST NOT*
be remembered for a different source, and *SHOULD NOT* be remembered across
sessions without a visible trust decision. A non-interactive host *MUST NOT*
evaluate remote TSL at all unless the operator has enabled it explicitly.

`encoding: "gpg"` protects confidentiality only: it is not an execution
boundary, provides no provenance, and hides the body from review. Consumers
*MUST NOT* automatically decrypt remote models, and *SHOULD* refuse encrypted
bodies when a restricted "safe" profile is requested.

**Includes.** An `#include` is remote code execution for GLSL/WGSL and plain
remote code execution for TSL. Therefore:

* A TSL body *MUST NOT* contain `#include`. Consumers *MUST* reject such a file
  rather than fetching the include.
* Consumers *MUST NOT* fetch remote includes for a model unless the user (or
  operator) has explicitly enabled network includes for that model. The default
  *MUST* be disabled.
* When enabled, consumers *MUST* use HTTPS and restrict requests and redirects
  to an explicit host allowlist. For `github.com/` includes, consumers *MUST*
  require an immutable 40-hex commit hash; branch names and tags are mutable
  and *MUST* be rejected. For hosts that serve unversioned content (such as
  `lygia.xyz`), consumers *MUST* verify a content hash when one is supplied,
  cache by content hash, and warn that the content is mutable. Consumers
  *MUST* enforce per-include and total byte limits plus an include-count limit
  and a request timeout, and *MUST* fail closed if an include cannot be fetched
  or verified. Silently substituting or omitting an include is forbidden.
* A consumer *MAY* offer an explicit, off-by-default override that accepts a
  mutable branch or tag reference. When it does, the consumer *MUST* warn for
  that specific source (naming the URL and, when available, a content hash)
  and the override *MUST NOT* apply to other models or to later sessions
  without a new decision. A restricted "safe" profile *MUST NOT* allow mutable
  includes at all.
* A future revision *SHOULD* define a self-contained bundle that embeds all
  included sources and their hashes so that rendering requires no network
  access.

**Resource limits.** Consumers *SHOULD* enforce limits before rendering:
maximum file size, maximum decoded body size, maximum include bytes and count,
maximum materials, finite `min`/`max` values with `min < max`, a maximum model
extent, and a maximum slice/sample count derived from the requested resolution.
A model *MUST NOT* be able to make a consumer allocate unbounded memory or run
unbounded work merely by declaring large bounds or by using a compression bomb
(`encoding: "gzip"`). Consumers *SHOULD* report which limit was exceeded.

**Untrusted header strings.** `materials`, `title`, `author`, `notes`,
`license`, and `options` are untrusted display data. Consumers *MUST NOT*
interpret them as HTML, shell code, format strings, or path fragments. In
particular: material names used in filenames *MUST* be reduced to a safe,
length-limited character set and *MUST NOT* contain path separators or `..`;
material names rendered in a user interface *MUST* be set as text, never as
markup; and header strings written to a terminal or log *MUST* have control
characters escaped.

**Provenance.** A future revision *SHOULD* define a `sha256` value over the
canonical file bytes and an optional detached signature, so that hosts can
show whether a model is signed and whether it has changed since it was
approved. In the absence of a signature, hosts *SHOULD* display a content hash
and the origin for any model loaded from a URL.

**Unknown values.** See the `language` key above: a consumer that does not
support `"tsl"` (or any unknown `language`) *MUST* reject the file with a clear
error and *MUST NOT* silently interpret it as another language.

**Versioning.** See `tslVersion` above.

## Determinism and slicing

The renderer modifies the parameter values of this function on each slice of the design in order
to calculate the amount of material needed at each point in 3D space. It is
free to “zoom in” to any portion of the design to get as much detail as
necessary to generate the model. This is why IRMF shaders have infinite
resolution. The renderer can get as much detail from the shader as it needs
in order to manufacture the part within the alloted timeframe. Higher
resolution models typically take more time to manufacture, so the same
IRMF shader can be used to create quick prototypes or highly-detailed,
final production-quality parts.

Because a TSL body is evaluated once to build a node graph (rather than once
per sample), a renderer *MUST NOT* rely on re-running the JavaScript for each
voxel; per-sample decisions are expressed in the graph itself.

## Revision history

* **1.0** — original specification (GLSL ES).
* **1.0 amendment** — added WGSL as a `language` value and documented the
  WGSL entry-point signatures.
* **1.0 amendment** — added TSL as a `language` value, with the `tslVersion`
  key, the TSL entry-point contract, and the execution-model, determinism,
  and security requirements above. `"irmf"` remains `"1.0"`, matching how
  WGSL was added.
* **1.0 amendment** — clarified the trust model: explicit per-source consent
  and sandboxing for TSL, no `#include` in TSL bodies, opt-in network includes
  with host allowlisting/pinning/limits (plus an optional explicit override
  for mutable refs), resource limits, and the rule that header strings are
  untrusted display data. Added rejection test vectors in
  [`test-vectors/`](test-vectors/).

## Canonical example

See [`examples/tsl-sphere.irmf`](examples/tsl-sphere.irmf) and
[`examples/tsl-two-material.irmf`](examples/tsl-two-material.irmf).
