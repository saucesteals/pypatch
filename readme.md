<div align="center">

<img src="./assets/logo.png" alt="logo" height="100px" />
<h1 align="center">PyPatch</h1>
<h3>Inject arbitrary code into any running python interpreter</h3>
</div>

---

## Overview

pypatch includes:

- [Python injector](https://github.com/saucesteals/pypatch/blob/main/pypatch) to run arbitrary code with any python3x.dll
- [Process injector library](https://github.com/saucesteals/pypatch/blob/main/inject) to inject dlls
- [Injector utility](https://github.com/saucesteals/pypatch/blob/main/cmd/injector) to simplify the entire process

## Usage

- Write the payload you want to inject in `inject.py`
- Compile a DLL with the payload with `make dll`
- Compile an injector (in `cmd/injector`) with your `program` and `dll` paths

## Contributing

Contributions are welcome!

- **[Submit Pull Requests](https://github.com/saucesteals/pypatch/pulls)**
- **[Report Issues](https://github.com/saucesteals/pypatch/issues)**

## License

Distributed under the GNU GPL v3.0 License. See `LICENSE` for more information.
