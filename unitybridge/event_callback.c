#include <stdint.h>

#include "event_callback.h"

#include "_cgo_export.h"

extern void goEventCallback(GoUint64 event_code, GoSlice data, GoUint64 tag);

void cEventCallback(uint64_t event_code, uintptr_t data, int length, uint64_t tag) {
    GoSlice data_slice;
    data_slice.data = (void *)data;
    data_slice.len = length;
    data_slice.cap = length;

    goEventCallback(event_code, data_slice, tag);
}
