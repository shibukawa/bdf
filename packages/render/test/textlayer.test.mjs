import { test } from "node:test";
import assert from "node:assert/strict";
import { joinRuns } from "../dist/textlayer.js";

// Sep values (docs/spec.md §7.9): NONE 0, SPACE 1, BREAK 2.
const run = (text, sep, layer = "p1") => ({ ordinal: 0, sep, text, layer });

test("joinRuns puts spaces at line breaks and newlines at paragraph breaks", () => {
  assert.equal(joinRuns([run("Hello", 2), run("world", 1), run("!", 0), run("Next", 2)]), "Hello world!\nNext");
});

test("joinRuns drops the first run's own separator", () => {
  assert.equal(joinRuns([run("a", 2), run("b", 2)]), "a\nb");
  assert.equal(joinRuns([run("a", 1)]), "a");
  assert.equal(joinRuns([]), "");
});

test("joinRuns separates pages with a blank line", () => {
  assert.equal(joinRuns([run("end of page", 0, "p1"), run("top of next", 2, "p2")]), "end of page\n\ntop of next");
});
