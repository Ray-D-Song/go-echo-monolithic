// SPA Router with URL monitoring like Vue/React Router
class SimpleRouter {
  constructor() {
    this.routes = {
      '/': 'home',
      '/home': 'home',
      '/about': 'about'
    };
    this.currentRoute = '/';
    this.init();
  }

  init() {
    // Handle browser back/forward buttons
    window.addEventListener('popstate', (e) => {
      this.handleRouteChange(window.location.pathname);
    });

    // Handle tab clicks
    document.addEventListener('click', (e) => {
      if (e.target.matches('[data-route]')) {
        e.preventDefault();
        const route = e.target.getAttribute('data-route');
        this.navigateTo(`/${route}`);
      }
    });

    // Load initial route
    this.handleRouteChange(window.location.pathname);
  }

  navigateTo(path) {
    if (this.currentRoute === path) return;
    
    // Update browser URL without reload
    window.history.pushState(null, '', path);
    this.handleRouteChange(path);
  }

  handleRouteChange(path) {
    // Normalize path
    const normalizedPath = path === '/' ? '/' : path;
    const tabName = this.routes[normalizedPath] || this.routes['/'];
    
    this.currentRoute = normalizedPath;
    this.switchTab(tabName);
  }

  switchTab(tabName) {
    // Update active tab
    document.querySelectorAll('.nav-tab').forEach(tab => {
      tab.classList.remove('active');
    });
    const activeTab = document.querySelector(`[data-route="${tabName}"]`);
    if (activeTab) {
      activeTab.classList.add('active');
    }

    // Update active page
    document.querySelectorAll('.page').forEach(page => {
      page.classList.remove('active');
    });
    const activePage = document.getElementById(tabName);
    if (activePage) {
      activePage.classList.add('active');
    }

    this.loadContent(tabName);
    
    // Log route change for debugging
    console.log(`Route changed to: ${this.currentRoute} -> ${tabName}`);
  }

  loadContent(tabName) {
    const homeContainer = document.getElementById('home');
    const aboutContainer = document.getElementById('about');

    if (tabName === 'home' && (!homeContainer.hasChildNodes() || homeContainer.innerHTML.includes('Loading'))) {
      this.renderHomePage(homeContainer);
    }

    if (tabName === 'about' && (!aboutContainer.hasChildNodes() || aboutContainer.innerHTML.includes('Loading'))) {
      this.renderAboutPage(aboutContainer);
    }
  }

  renderHomePage(container) {
    container.innerHTML = `
      <h2>Overview</h2>
      
      <div class="content-section">
        <h3>Features</h3>
        <ul>
          <li>Tab-based navigation</li>
          <li>Minimal JavaScript</li>
          <li>Clean interface</li>
          <li>Static file serving</li>
        </ul>
      </div>

      <div class="content-section">
        <h3>Technology</h3>
        <div class="tech-list">
          <span class="tech-item">Go</span>
          <span class="tech-item">Echo</span>
          <span class="tech-item">GORM</span>
          <span class="tech-item">JavaScript</span>
          <span class="tech-item">HTML</span>
          <span class="tech-item">CSS</span>
        </div>
      </div>

      <div class="content-section">
        <h3>Purpose</h3>
        <p>This is a simple single page application for testing Echo framework's static file serving capabilities.</p>
      </div>
    `;
  }

  renderAboutPage(container) {
    container.innerHTML = `
      <h2>Implementation</h2>
      
      <div class="content-section">
        <h3>Structure</h3>
        <p>Simple tab-based navigation without URL routing. Content is loaded dynamically using JavaScript.</p>
      </div>

      <div class="content-section">
        <h3>Usage</h3>
        <p>Serve the web directory as static files through Echo. The application handles client-side navigation internally.</p>
      </div>

      <div class="content-section">
        <h3>Files</h3>
        <ul>
          <li>index.html - Main HTML structure</li>
          <li>spa-example.js - Router and page logic</li>
        </ul>
      </div>
    `;
  }
}

// Initialize the router when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  new SimpleRouter();
});