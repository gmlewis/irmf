# IRMF security rejection test vectors

These files are intentionally hostile or malformed. They exist so that every
consumer (editor, slicers, linters, converters) can be tested against the same
expectations. Consumers **must fail closed**; they must never fall back to a
different interpretation, fetch code, or render partial output.

| File | Expectation |
| --- | --- |
| `reject-unknown-language.irmf` | Reject with an "unsupported language" error. Never treat the body as GLSL. |
| `reject-tsl-include.irmf` | Reject before evaluating JavaScript or fetching anything. A TSL body must not contain `#include`. |
| `reject-unpinned-include.irmf` | Reject because the include targets a branch, not an immutable commit. Must not fetch it when includes are disabled (the default). With includes enabled, only an explicit, off-by-default override that warns for this specific source may accept it; a restricted "safe" profile must reject it. |
| `reject-huge-bounds.irmf` | Reject because the declared extent exceeds the consumer's model-extent limit. Must not attempt to allocate/slice the implied volume. |
| `reject-gpg-untrusted.irmf` | Refuse to auto-decrypt or evaluate the payload. Encrypted bodies are not a trust boundary; only decrypt after an explicit trust decision, and never in a restricted "safe" profile. |
| `reject-bad-material-path.irmf` | If accepted, any output filename derived from the material name must stay inside the output directory and must not contain `..` or path separators. Prefer rejecting or sanitizing the name. |
| `display-material-name.irmf` | Load and render normally, but display the material name as literal text. It must never be inserted as HTML; no script/event handler may run. |

Notes:

* A consumer that does not support a feature (for example GPG or remote
  includes) must reject the file clearly rather than ignoring the feature.
* Network-access test vectors must also be run with networking disabled to
  verify that no request is attempted.
* `reject-gpg-untrusted.irmf` contains a dummy PGP armored block; it is not
  decryptable and is only present to exercise the policy.
