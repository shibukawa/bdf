import type { DCValue } from "./types.js";

/** The values of a Dublin Core element as an array (empty when the element is absent). */
export function dcValues(v: DCValue | undefined): string[] {
  if (v === undefined) return [];
  return typeof v === "string" ? [v] : v;
}
