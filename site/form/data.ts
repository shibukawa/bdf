/** A value of JSON. */
export type Json = null | boolean | number | string | Json[] | JsonObject;
export type JsonObject = { [name: string]: Json };

/** Template field names are data, including __proto__ and constructor. */
export function valueOf(object: JsonObject, name: string): Json | undefined {
  return Object.hasOwn(object, name) ? object[name] : undefined;
}

export function setValue(object: JsonObject, name: string, value: Json) {
  Object.defineProperty(object, name, { value, writable: true, enumerable: true, configurable: true });
}

/** Walk the data once, instead of searching its entire JSON for every picture. */
export function pictureFiles(data: Json, pictures: ReadonlyMap<string, Uint8Array>): Record<string, Uint8Array> {
  if (!pictures.size) return {};
  const found = new Map<string, Uint8Array>();
  const pending: Json[] = [data];
  while (pending.length) {
    const value = pending.pop()!;
    if (typeof value === "string") {
      const bytes = pictures.get(value);
      if (bytes) found.set(value, bytes);
    } else if (value !== null && typeof value === "object") {
      for (const child of Object.values(value)) pending.push(child);
    }
  }
  return Object.fromEntries(found);
}
