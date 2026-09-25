// Encode-only glue around libwebp for the wasm module translated into
// imgconv/internal/webpw. Input is RGBA (4 bytes per pixel, row stride w*4).
// Returns a malloc'd buffer (free it with free) and writes its length to *size.
#include <stdlib.h>
#include "webp/encode.h"

__attribute__((export_name("encode")))
uint8_t *encode(const uint8_t *rgba, int w, int h, size_t *size, int quality, int method, int lossless, int exact) {
    *size = 0;
    WebPConfig config;
    if (!WebPConfigInit(&config)) return NULL;
    config.quality = (float)quality;
    config.method = method;
    config.lossless = lossless;
    config.exact = exact;
    if (!WebPValidateConfig(&config)) return NULL;

    WebPPicture picture;
    if (!WebPPictureInit(&picture)) return NULL;
    picture.width = w;
    picture.height = h;
    picture.use_argb = 1;
    if (!WebPPictureImportRGBA(&picture, rgba, w * 4)) {
        WebPPictureFree(&picture);
        return NULL;
    }
    WebPMemoryWriter writer;
    WebPMemoryWriterInit(&writer);
    picture.writer = WebPMemoryWrite;
    picture.custom_ptr = &writer;
    int ok = WebPEncode(&config, &picture);
    WebPPictureFree(&picture);
    if (!ok) {
        WebPMemoryWriterClear(&writer);
        return NULL;
    }
    *size = writer.size;
    return writer.mem;
}
