#ifndef FUNCTION_CALLER_H_
#define FUNCTION_CALLER_H_

#include <stdbool.h>
#include <stdint.h>

#include "event_callback.h"

// Wraper to call C function pointers as CGO can not do it by itself
// apparently.

void CreateUnityBridgeCaller(void *f, const char *name, bool debuggable,
                             const char *log_path);

void DestroyUnityBridgeCaller(void *f);

bool UnityBridgeInitializeCaller(void *f);

void UnityBridgeUninitializeCaller(void *f);

void UnitySendEventCaller(void *f, uint64_t event_code, uintptr_t data,
                          int length, uint64_t tag);

void UnitySendEventWithStringCaller(void *f, uint64_t event_code,
                                    const char *data, uint64_t tag);

void UnitySendEventWithNumberCaller(void *f, uint64_t event_code, uint64_t data,
                                    uint64_t tag);

void UnitySetEventCallbackCaller(void *f, uint64_t event_code,
                                 EventCallback event_callback);

char* UnityGetSecurityKeyByKeyChainIndexCaller(void *f, int index);

#endif  // FUNCTION_CALLER_H_