#include "function_caller.h"

void CreateUnityBridgeCaller(void *f, const char *name, bool debuggable,
                             const char *log_path) {
  ((void (*)(const char *, bool, const char *))f)(name, debuggable, log_path);
}

void DestroyUnityBridgeCaller(void *f) { ((void (*)())f)(); }

bool UnityBridgeInitializeCaller(void *f) { return ((bool (*)())f)(); }

void UnityBridgeUninitializeCaller(void *f) { ((void (*)())f)(); }

void UnitySendEventCaller(void *f, uint64_t event_code, uintptr_t data,
                          int length, uint64_t tag) {
  ((void (*)(uint64_t, uintptr_t, int, uint64_t))f)(event_code, data, length,
                                                    tag);
}

void UnitySendEventWithStringCaller(void *f, uint64_t event_code,
                                    const char *data, uint64_t tag) {
  ((void (*)(uint64_t, const char *, uint64_t))f)(event_code, data, tag);
}

void UnitySendEventWithNumberCaller(void *f, uint64_t event_code, uint64_t data,
                                    uint64_t tag) {
  ((void (*)(uint64_t, uint64_t, uint64_t))f)(event_code, data, tag);
}

void UnitySetEventCallbackCaller(void *f, uint64_t event_code,
                                 EventCallback event_callback) {
  ((void (*)(uint64_t, EventCallback))f)(event_code, event_callback);
}

char* UnityGetSecurityKeyByKeyChainIndexCaller(void *f, int index) {
  return (char*)((uintptr_t(*)(int))f)(index);
}
