#include "dll.h"

extern void OnProcessAttach();

static DWORD WINAPI OnAttachThread(LPVOID lpParam) {
  OnProcessAttach();
  return 0;
}

BOOL WINAPI DllMain(HINSTANCE hinstDLL, DWORD fdwReason, LPVOID lpReserved) {
  switch (fdwReason) {
  case DLL_PROCESS_ATTACH:
    CreateThread(NULL, 0, OnAttachThread, NULL, 0, NULL);
    break;
  }
  return TRUE;
}