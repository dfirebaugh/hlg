const template = document.createElement('template');
template.innerHTML = `
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css"
integrity="sha512-iecdLmaskl7CVkqkXNQ/ZH/XLlvWZOJyj7Yy7tcenmpD1ypASozpmT/E0iPtmFIB46ZmdtAc9eNBvH0H/ZpiBw=="
crossorigin="anonymous" referrerpolicy="no-referrer" />

<style>
    .top-bar {
        display: grid;
        align-items: center;
        justify-content: space-between;
        grid-template-columns: 1fr 4fr 1fr;
        background-color: var(--secondary-background-color) !important;
    }

    a {
        color: var(--primary-color);
    }

    .icon-button {
        border: none;
        background: none;
        padding: 0;
        margin: 0;
        outline: none;
        cursor: pointer;
        display: inline-flex;
        align-items: center;
        justify-content: center;
    }

    .icon-button i.fa {
        font-size: 1.2em;
        color: var(--primary-color);
    }

    .icon-button:hover, .icon-button:active {
        background-color: rgba(0, 0, 0, 0.1);
    }

    #site-title:hover {
        cursor: pointer;
        text-decoration: underline;
    }

    .menu-title {
        display:flex;
        justify-content: center;
    }
    left-buttons {
        padding-left: 2rem;
    }

    .right-buttons {
        justify-self: end;
        padding-right: 2rem;
    }
</style>
    <span class="top-bar">
        <left-buttons>
            <button id="sidebar-toggle" onclick="hamburgerClick()" class="icon-button" type="button" title="Toggle Table of Contents"
                aria-label="Toggle Table of Contents" aria-controls="sidebar" aria-expanded="true">
                <i class="fa fa-bars"></i>
            </button>
        </left-buttons>

        <h1 class="menu-title  md-6">
            <a id="site-title">
            </a>
        </h1>

        <div class="right-buttons">
            <a id="repo-link" href="" target="_blank" class="icon-button" title="GitHub Repository">
                <i id="repo-icon" src="" alt="GitHub" />
            </a>
            <a id="edit-page" href="" target="_blank" class="icon-button" title="Edit Page">
                <i class="fa fa-edit"></i>
            </a>
        </div>
    </span>
`;

class TopBar extends HTMLElement {
    static get observedAttributes() {
        return [
            'title',
            'url',
            'theme',
            'hideSideBar',
            'repoIcon',
            'repoURL',
            'editURL',
        ];
    }

    constructor() {
        super();
        this.attachShadow({ mode: 'open' });
        this.shadowRoot.appendChild(template.content.cloneNode(true));
    }

    connectedCallback() {
        this.updateTitleAndUrl();
        this.updateTheme();
        this.hideSideBar();
        this.updateRepoAndEditLinks();
    }

    attributeChangedCallback(name, oldValue, newValue) {
        if (oldValue !== newValue) {
            if (name === 'title' || name === 'url') {
                this.updateTitleAndUrl();
            } else if (name === 'theme') {
                this.updateTheme();
            } else if (name === 'repoIcon' || name === 'repoURL' || name === 'editURL') {
                this.updateRepoAndEditLinks();
            }
        }
    }

    updateTitleAndUrl() {
        const titleElement = this.shadowRoot.querySelector("#site-title");
        const title = this.getAttribute('title');
        const url = this.getAttribute('url');

        if (title) titleElement.innerText = title;
        if (url) titleElement.href = url;
    }

    hideSideBar() {
        if (this.getAttribute('hideSideBar') !== 'true') return;

        const sidebarToggle = this.shadowRoot.querySelector("#sidebar-toggle");
        if (sidebarToggle) {
            sidebarToggle.style.display = 'none';
        }
    }

    updateTheme() {
        const theme = this.getAttribute('theme');
        if (theme) this.shadowRoot.host.className = theme;
    }

    updateRepoAndEditLinks() {
        const repoLink = this.shadowRoot.querySelector("#repo-link");
        const repoIcon = this.shadowRoot.querySelector("#repo-icon");
        const editLink = this.shadowRoot.querySelector("#edit-page");

        const repoURL = this.getAttribute('repoURL');
        const repoIconClass = this.getAttribute('repoIcon');
        const editURL = this.getAttribute('editURL');

        if (repoURL) {
            repoLink.href = repoURL;
            repoLink.style.display = 'inline';
            if (repoIconClass) {
                repoIcon.className = `fa-brands ${repoIconClass}`;
                repoIcon.style.display = 'inline';
            } else {
                repoIcon.style.display = 'none';
            }
        } else {
            repoLink.style.display = 'none';
        }

        if (editURL) {
            editLink.href = editURL;
            editLink.style.display = 'inline';
        } else {
            editLink.style.display = 'none';
        }
    }
}

window.customElements.define('bb-topbar', TopBar);
