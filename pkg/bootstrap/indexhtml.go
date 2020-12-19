package bootstrap

const indexHtml = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Go wasm</title>
    <link rel="stylesheet" href="style.css"/>
    <script src="_wasm_exec.js"></script>
    <script>
        // Start webassembly
        if (!WebAssembly.instantiateStreaming) { // polyfill
            WebAssembly.instantiateStreaming = async (resp, importObject) => {
                const source = await (await resp).arrayBuffer();
                return await WebAssembly.instantiate(source, importObject);
            };
        }
        const go = new Go();
        let mod, inst;
        WebAssembly.instantiateStreaming(fetch("_binary"), go.importObject).then(async (result) => {
            mod = result.module;
            inst = result.instance;
            await go.run(inst);
        }).catch((err) => {
            console.error(err);
        });
    </script>
</head>
    <body></body>
</html>`