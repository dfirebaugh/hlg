const template = document.createElement('template');
template.innerHTML = `
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css"
integrity="sha512-iecdLmaskl7CVkqkXNQ/ZH/XLlvWZOJyj7Yy7tcenmpD1ypASozpmT/E0iPtmFIB46ZmdtAc9eNBvH0H/ZpiBw=="
crossorigin="anonymous" referrerpolicy="no-referrer" />

<style>
button {
    background: none;
    border: none;
    color: var(--primary-color);
}
</style>

<aside id="sidenav" class="">
    <div class="">
        <button id="sidebar-close-sidebar" onclick="hamburgerClick()" class="icon-button" type="button"
            title="Toggle Table of Contents" aria-label="Toggle Table of Contents"
            aria-controls="sidebar" aria-expanded="true">
            <i class="fa fa-xmark"></i>
        </button>
        <ol class="navbar prose overflow-hidden" style="padding-left: 0;">
            <slot name="nav-content"></slot>
        </ol>
    </div>
</aside>`;

class NavMenu extends HTMLElement {
    constructor() {
        super();
        this.attachShadow({ mode: "open" });
        this.shadowRoot.appendChild(template.content.cloneNode(true));
    }
}
window.customElements.define('bb-navmenu', NavMenu);
