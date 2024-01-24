
const template = document.createElement('template');
template.innerHTML = `
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css"
integrity="sha512-iecdLmaskl7CVkqkXNQ/ZH/XLlvWZOJyj7Yy7tcenmpD1ypASozpmT/E0iPtmFIB46ZmdtAc9eNBvH0H/ZpiBw=="
crossorigin="anonymous" referrerpolicy="no-referrer" />

<style>
a {
    position:fixed;
    height: 100vh;
    display: flex;
    align-items: center;
    width: 5vw;
    justify-content: center;
}

a:hover {
    background-color: var(--secondary-background-color);
    cursor: pointer;
}

a {
    color: var(--primary-color);
}

.nav-wide-wrapper {
    display: flex;
}
</style>

<span class="nav-wide-wrapper" aria-label="Page navigation">
    <a rel="prev" class="nav-chapters previous" title="Previous chapter" aria-label="Previous chapter"
        aria-keyshortcuts="Left">
        <i class="fa fa-angle-left"></i>
    </a>
    <a rel="next" class="nav-chapters next" title="Next chapter" aria-label="Next chapter"
        aria-keyshortcuts="Right">
        <i class="fa fa-angle-right"></i>
    </a>
</span>`;

class NavBtns extends HTMLElement {
    static get observedAttributes() {
        return ['nav-type', 'href', 'theme'];
    }

    constructor() {
        super();
        this.attachShadow({ mode: 'open' }).appendChild(template.content.cloneNode(true));
    }

    connectedCallback() {
        this.updateLinks();
    }

    attributeChangedCallback(name, oldValue, newValue) {
        if (oldValue !== newValue) {
            this.updateLinks();
        }
    }

    updateLinks() {
        const navType = this.getAttribute('nav-type');
        const href = this.getAttribute('href');
        const prev = this.shadowRoot.querySelector('.previous');
        const next = this.shadowRoot.querySelector('.next');
        const wrapper = this.shadowRoot.querySelector('.nav-wide-wrapper');


        if (navType === 'prev') {
            next?.remove();
            wrapper.style.justifyContent = 'left';
            if (prev) prev.href = href;
        } else if (navType === 'next') {
            prev?.remove();
            wrapper.style.justifyContent = 'right';
            if (next) next.href = href;
        }
    }
}

window.customElements.define('bb-navbtn', NavBtns);
