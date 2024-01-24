const template = document.createElement('template');
template.innerHTML = `
<iframe id='wasm-iframe' 
    allow="autoplay"
    scrolling="no"
>
</iframe>
`;

class GoWasmCanvas extends HTMLElement {
    constructor() {
        super();
        this.attachShadow({ mode: 'open' });
        this.shadowRoot.appendChild(template.content.cloneNode(true));
        this.shadowRoot.getElementById('wasm-iframe').width = this.getAttribute('width');
        this.shadowRoot.getElementById('wasm-iframe').height = this.getAttribute('height');
    }

    connectedCallback() {
        const doc = this.shadowRoot.getElementById('wasm-iframe').contentWindow.document
        doc.open();
        doc.write(this.buildInnerIFrame(this.getAttribute('src')))
        doc.close()
    }

    buildInnerIFrame(src) {
        return `
<!DOCTYPE html>
        <script src="js/vendor/wasm_exec.js"></script>

        <script>
        // Polyfill
        if (!WebAssembly.instantiateStreaming) {
            WebAssembly.instantiateStreaming = async (resp, importObject) => {
                const source = await (await resp).arrayBuffer();
                return await WebAssembly.instantiate(source, importObject);
            };
        }

        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("${src}"), go.importObject).then(result => {
            go.run(result.instance);
        });
        </script>
  `;
    }
}

window.customElements.define('go-wasm-canvas', GoWasmCanvas);
