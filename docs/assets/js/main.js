// Theme Toggle Implementation
function initThemeToggle() {
    const themeToggle = document.getElementById('themeToggle');
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)');

    // Set initial theme based on user preference
    function setInitialTheme() {
        if (localStorage.getItem('theme') === 'dark' || 
            (!localStorage.getItem('theme') && prefersDark.matches)) {
            document.body.classList.add('dark-theme');
            themeToggle.querySelector('i').classList.replace('fa-moon', 'fa-sun');
        }
    }

    // Toggle theme
    function toggleTheme() {
        document.body.classList.toggle('dark-theme');
        const icon = themeToggle.querySelector('i');
        
        if (document.body.classList.contains('dark-theme')) {
            icon.classList.replace('fa-moon', 'fa-sun');
            localStorage.setItem('theme', 'dark');
        } else {
            icon.classList.replace('fa-sun', 'fa-moon');
            localStorage.setItem('theme', 'light');
        }
    }

    setInitialTheme();
    themeToggle.addEventListener('click', toggleTheme);
    prefersDark.addEventListener('change', setInitialTheme);
}

// Tab Functionality Implementation
function initTabFunctionality() {
    const tabButtons = document.querySelectorAll('.tab-btn');
    
    function switchTab(tabId) {
        // Update button states
        tabButtons.forEach(btn => {
            btn.classList.toggle('active', btn.getAttribute('data-tab') === tabId);
        });

        // Update content visibility
        document.querySelectorAll('.tab-content').forEach(content => {
            content.style.display = content.id === tabId ? 'block' : 'none';
        });
    }

    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            switchTab(button.getAttribute('data-tab'));
        });
    });
}

// Modal Functionality Implementation
function initModalFunctionality() {
    const modalOverlay = document.getElementById('modalOverlay');
    const modalMarkdown = document.getElementById('modalMarkdown');
    let currentScroll = 0;

    window.openModal = function(mdFile) {
        currentScroll = window.scrollY;
        modalMarkdown.setAttribute('src', mdFile);
        modalOverlay.style.display = 'flex';
        document.body.style.position = 'fixed';
        document.body.style.top = `-${currentScroll}px`;
        document.body.style.width = '100%';

    };

    window.closeModal = function() {
        modalOverlay.style.display = 'none';
        modalMarkdown.removeAttribute('src');
        document.body.style.position = '';
        document.body.style.top = '';
        document.body.style.width = '';
        window.scrollTo(0, currentScroll);
    };

    // Close modal on escape key
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && modalOverlay.style.display === 'flex') {
            window.closeModal();
        }
    });

    // Close modal when clicking outside
    modalOverlay.addEventListener('click', (e) => {
        if (e.target === modalOverlay) {
            window.closeModal();
        }
    });
}

// Smooth Scroll Implementation
function initSmoothScroll() {
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function(e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            
            if (target) {
                const navHeight = document.querySelector('.navbar').offsetHeight;
                const targetPosition = target.getBoundingClientRect().top + window.pageYOffset;
                
                window.scrollTo({
                    top: targetPosition - navHeight,
                    behavior: 'smooth'
                });
            }
        });
    });
}

// Scroll Spy Implementation
function initScrollSpy() {
    const navLinks = document.querySelectorAll('.nav-link[href^="#"]');
    const sections = document.querySelectorAll('section[id]');
    
    function updateActiveLink() {
        const navHeight = document.querySelector('.navbar').offsetHeight;
        const fromTop = window.scrollY + navHeight;

        sections.forEach(section => {
            const sectionTop = section.offsetTop - navHeight;
            const sectionBottom = sectionTop + section.offsetHeight;
            const id = section.getAttribute('id');
            const correspondingLink = document.querySelector(`.nav-link[href="#${id}"]`);

            if (fromTop >= sectionTop && fromTop < sectionBottom && correspondingLink) {
                navLinks.forEach(link => link.classList.remove('active'));
                correspondingLink.classList.add('active');
            }
        });

        // Special case for page bottom
        const pageBottom = window.scrollY + window.innerHeight === document.documentElement.scrollHeight;
        if (pageBottom) {
            navLinks.forEach(link => link.classList.remove('active'));
            navLinks[navLinks.length - 1].classList.add('active');
        }
    }

    // Throttle scroll event
    let ticking = false;
    window.addEventListener('scroll', () => {
        if (!ticking) {
            window.requestAnimationFrame(() => {
                updateActiveLink();
                ticking = false;
            });
            ticking = true;
        }
    });

    // Update active link on page load
    updateActiveLink();
}

// Handle copy functionality for code blocks
function initCodeCopy() {
    document.querySelectorAll('pre code').forEach((block) => {
        const copyButton = document.createElement('button');
        copyButton.className = 'copy-button';
        copyButton.innerHTML = '<i class="fas fa-copy"></i>';
        
        copyButton.addEventListener('click', () => {
            navigator.clipboard.writeText(block.textContent).then(() => {
                copyButton.innerHTML = '<i class="fas fa-check"></i>';
                setTimeout(() => {
                    copyButton.innerHTML = '<i class="fas fa-copy"></i>';
                }, 2000);
            }).catch((err) => {
                console.error('Failed to copy text: ', err);
                copyButton.innerHTML = '<i class="fas fa-exclamation-triangle"></i>';
                setTimeout(() => {
                    copyButton.innerHTML = '<i class="fas fa-copy"></i>';
                }, 2000);
            });
        });

        const wrapper = document.createElement('div');
        wrapper.className = 'code-wrapper';
        block.parentNode.insertBefore(wrapper, block);
        wrapper.appendChild(block);
        wrapper.appendChild(copyButton);
    });
}

// window.ZeroMD = {
//     config: {
//         cssUrls: [
//             'https://cdn.jsdelivr.net/gh/sindresorhus/github-markdown-css/github-markdown.css'
//         ],
//         markedOptions: {
//             breaks: true,
//             gfm: true
//         }
//     }
// };

// Initialize code copy functionality when zero-md components are loaded


document.addEventListener('zero-md-rendered', initCodeCopy);

// Wait for DOM to be fully loaded
document.addEventListener('DOMContentLoaded', function() {
    // Initialize AOS (Animate On Scroll)
    AOS.init({
        duration: 800,
        once: true,
        offset: 50
    });

    initThemeToggle();
    initTabFunctionality();
    initModalFunctionality();
    initSmoothScroll();
    initScrollSpy();
});

