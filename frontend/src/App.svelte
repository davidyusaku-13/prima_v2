<script>
  import { onMount } from 'svelte';

  let email = '';
  let subscribed = false;
  let loading = false;
  let isMenuOpen = false;
  let visibleSections = new Set();

  function toggleMenu() {
    isMenuOpen = !isMenuOpen;
  }

  function closeMenu() {
    isMenuOpen = false;
  }

  // Track event in GA4
  function trackEvent(eventName, params = {}) {
    if (typeof window.gtag === 'function') {
      window.gtag('event', eventName, params);
    }
  }

  async function handleSubscribe() {
    if (!email.trim()) return;

    loading = true;
    // Track subscription attempt
    trackEvent('begin_checkout', {
      currency: 'USD',
      value: 0,
      items: [{ item_id: 'newsletter', item_name: 'Newsletter' }]
    });

    // Simulate subscription
    await new Promise(resolve => setTimeout(resolve, 1000));
    subscribed = true;
    loading = false;

    // Track successful subscription
    trackEvent('subscribe', {
      method: 'email',
      content_type: 'newsletter'
    });
  }

  function handleCTAClick(ctaName) {
    trackEvent('select_content', {
      content_type: 'button',
      item_id: ctaName
    });
  }

  onMount(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            visibleSections.add(entry.target.id);
            visibleSections = visibleSections;
          }
        });
      },
      { threshold: 0.1, rootMargin: '0px 0px -50px 0px' }
    );

    document.querySelectorAll('section[id]').forEach((section) => {
      observer.observe(section);
    });

    return () => observer.disconnect();
  });
</script>

<main class="min-h-screen bg-slate-50">
  <!-- Navigation -->
  <nav class="sticky top-0 z-40 bg-white/80 backdrop-blur-md border-b border-slate-200">
    <div class="max-w-6xl mx-auto px-6 py-4">
      <div class="flex items-center justify-between">
        <a href="#" class="text-2xl font-bold bg-gradient-to-r from-slate-900 to-slate-700 bg-clip-text text-transparent">
          Prima
        </a>

        <!-- Desktop Menu -->
        <div class="hidden md:flex items-center gap-8">
          <a href="#features" class="text-slate-600 hover:text-slate-900 font-medium transition-colors duration-200">Features</a>
          <a href="#about" class="text-slate-600 hover:text-slate-900 font-medium transition-colors duration-200">About</a>
          <a href="#contact" class="text-slate-600 hover:text-slate-900 font-medium transition-colors duration-200">Contact</a>
        </div>

        <!-- Desktop CTA -->
        <div class="hidden md:block">
          <button
            onclick={() => handleCTAClick('get_started_nav')}
            class="px-6 py-2.5 bg-slate-900 text-white font-semibold rounded-xl hover:bg-slate-800 hover:shadow-lg hover:-translate-y-0.5 transition-all duration-200"
          >
            Get Started
          </button>
        </div>

        <!-- Mobile Menu Button -->
        <button
          onclick={toggleMenu}
          class="md:hidden p-2.5 text-slate-600 hover:text-slate-900 hover:bg-slate-100 rounded-xl transition-all duration-200"
          aria-label="Toggle menu"
        >
          {#if isMenuOpen}
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          {:else}
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          {/if}
        </button>
      </div>
    </div>
  </nav>

  <!-- Mobile Menu Overlay -->
  {#if isMenuOpen}
    <div class="fixed inset-0 z-50 md:hidden">
      <div
        class="absolute inset-0 bg-slate-900/40 backdrop-blur-sm"
        onclick={closeMenu}
        role="button"
        tabindex="-1"
        onkeydown={(e) => e.key === 'Escape' && closeMenu()}
      ></div>

      <div class="absolute right-0 top-0 h-full w-80 bg-white shadow-2xl">
        <div class="p-6">
          <div class="flex items-center justify-between mb-8">
            <span class="text-xl font-bold text-slate-900">Menu</span>
            <button
              onclick={closeMenu}
              class="p-2.5 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-xl transition-all duration-200"
              aria-label="Close menu"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div class="flex flex-col gap-2">
            <a
              href="#features"
              onclick={closeMenu}
              class="px-4 py-3 text-slate-600 hover:text-slate-900 hover:bg-slate-50 font-medium rounded-xl transition-all duration-200"
            >Features</a>
            <a
              href="#about"
              onclick={closeMenu}
              class="px-4 py-3 text-slate-600 hover:text-slate-900 hover:bg-slate-50 font-medium rounded-xl transition-all duration-200"
            >About</a>
            <a
              href="#contact"
              onclick={closeMenu}
              class="px-4 py-3 text-slate-600 hover:text-slate-900 hover:bg-slate-50 font-medium rounded-xl transition-all duration-200"
            >Contact</a>
          </div>

          <div class="mt-8 pt-6 border-t border-slate-100">
            <button
              onclick={() => { closeMenu(); handleCTAClick('get_started_mobile'); }}
              class="w-full px-5 py-3.5 bg-slate-900 text-white font-semibold rounded-xl hover:bg-slate-800 hover:shadow-lg hover:-translate-y-0.5 transition-all duration-200"
            >
              Get Started
            </button>
          </div>
        </div>
      </div>
    </div>
  {/if}

  <!-- Hero Section -->
  <section class="px-6 py-24 md:py-32 max-w-6xl mx-auto">
    <div class="max-w-3xl">
      <div class="inline-flex items-center gap-2 px-4 py-2 bg-blue-50 text-blue-700 rounded-full text-sm font-semibold mb-8 animate-fade-in">
        <span class="w-2 h-2 bg-blue-500 rounded-full animate-pulse"></span>
        Now in public beta
      </div>
      <h1 class="text-5xl md:text-7xl font-bold text-slate-900 leading-[1.1] mb-6 tracking-tight animate-fade-in-up">
        Build better products,
        <span class="bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">faster.</span>
      </h1>
      <p class="text-xl md:text-2xl text-slate-600 mb-10 leading-relaxed max-w-2xl animate-fade-in-up delay-100">
        Prima helps teams collaborate seamlessly and ship quality products.
        Streamline your workflow with powerful tools designed for modern development.
      </p>
      <div class="flex flex-wrap gap-4 animate-fade-in-up delay-200">
        <button
          onclick={() => handleCTAClick('start_free_trial')}
          class="px-8 py-4 bg-blue-600 text-white font-semibold rounded-xl hover:bg-blue-700 hover:shadow-lg hover:-translate-y-0.5 transition-all duration-200"
        >
          Start Free Trial
        </button>
        <button
          onclick={() => handleCTAClick('watch_demo')}
          class="px-8 py-4 border-2 border-slate-200 text-slate-700 font-semibold rounded-xl hover:border-slate-300 hover:bg-slate-50 transition-all duration-200"
        >
          Watch Demo
        </button>
      </div>
    </div>
  </section>

  <!-- Stats Section -->
  <section class="bg-white py-20 border-y border-slate-200">
    <div class="max-w-6xl mx-auto px-6">
      <div class="grid grid-cols-2 md:grid-cols-4 gap-8 md:gap-12">
        <div class="text-center group">
          <div class="text-4xl md:text-5xl font-bold text-slate-900 mb-2 group-hover:scale-110 transition-transform duration-300">10k+</div>
          <div class="text-slate-500 font-medium">Active Users</div>
        </div>
        <div class="text-center group">
          <div class="text-4xl md:text-5xl font-bold text-slate-900 mb-2 group-hover:scale-110 transition-transform duration-300">99.9%</div>
          <div class="text-slate-500 font-medium">Uptime</div>
        </div>
        <div class="text-center group">
          <div class="text-4xl md:text-5xl font-bold text-slate-900 mb-2 group-hover:scale-110 transition-transform duration-300">50M+</div>
          <div class="text-slate-500 font-medium">API Requests</div>
        </div>
        <div class="text-center group">
          <div class="text-4xl md:text-5xl font-bold text-slate-900 mb-2 group-hover:scale-110 transition-transform duration-300">24/7</div>
          <div class="text-slate-500 font-medium">Support</div>
        </div>
      </div>
    </div>
  </section>

  <!-- Features Section -->
  <section id="features" class="py-24 px-6 max-w-6xl mx-auto">
    <div class="text-center mb-16">
      <h2 class="text-4xl md:text-5xl font-bold text-slate-900 mb-4 tracking-tight">Everything you need</h2>
      <p class="text-xl text-slate-600 max-w-2xl mx-auto">Powerful features designed to help your team succeed and scale effortlessly.</p>
    </div>
    <div class="grid md:grid-cols-3 gap-8">
      <!-- Feature 1 -->
      <div class="group p-8 bg-white rounded-2xl border border-slate-200 hover:border-blue-200 hover:shadow-xl hover:shadow-blue-500/5 transition-all duration-300 hover:-translate-y-1">
        <div class="w-14 h-14 bg-blue-100 rounded-2xl flex items-center justify-center mb-6 group-hover:bg-blue-600 transition-colors duration-300">
          <svg class="w-7 h-7 text-blue-600 group-hover:text-white transition-colors duration-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
        </div>
        <h3 class="text-xl font-bold text-slate-900 mb-3">Lightning Fast</h3>
        <p class="text-slate-600 leading-relaxed">Built for speed with optimized performance that keeps your team productive and in the flow.</p>
      </div>

      <!-- Feature 2 -->
      <div class="group p-8 bg-white rounded-2xl border border-slate-200 hover:border-green-200 hover:shadow-xl hover:shadow-green-500/5 transition-all duration-300 hover:-translate-y-1">
        <div class="w-14 h-14 bg-green-100 rounded-2xl flex items-center justify-center mb-6 group-hover:bg-green-600 transition-colors duration-300">
          <svg class="w-7 h-7 text-green-600 group-hover:text-white transition-colors duration-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
          </svg>
        </div>
        <h3 class="text-xl font-bold text-slate-900 mb-3">Secure by Default</h3>
        <p class="text-slate-600 leading-relaxed">Enterprise-grade security with end-to-end encryption and SOC2 compliance built right in.</p>
      </div>

      <!-- Feature 3 -->
      <div class="group p-8 bg-white rounded-2xl border border-slate-200 hover:border-purple-200 hover:shadow-xl hover:shadow-purple-500/5 transition-all duration-300 hover:-translate-y-1">
        <div class="w-14 h-14 bg-purple-100 rounded-2xl flex items-center justify-center mb-6 group-hover:bg-purple-600 transition-colors duration-300">
          <svg class="w-7 h-7 text-purple-600 group-hover:text-white transition-colors duration-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
        </div>
        <h3 class="text-xl font-bold text-slate-900 mb-3">Team Collaboration</h3>
        <p class="text-slate-600 leading-relaxed">Work together seamlessly with real-time updates, comments, and powerful sharing features.</p>
      </div>
    </div>
  </section>

  <!-- About Section -->
  <section id="about" class="py-24 px-6 max-w-6xl mx-auto">
    <div class="grid md:grid-cols-2 gap-16 items-center">
      <div>
        <h2 class="text-4xl md:text-5xl font-bold text-slate-900 mb-6 tracking-tight">About Prima</h2>
        <p class="text-slate-600 mb-6 text-lg leading-relaxed">
          Prima was founded with a simple mission: to make software development
          more efficient and enjoyable for teams of all sizes.
        </p>
        <p class="text-slate-600 mb-6 text-lg leading-relaxed">
          Today, thousands of companies trust Prima to power their development
          workflows, from innovative startups to Fortune 500 enterprises.
        </p>
        <p class="text-slate-600 text-lg leading-relaxed">
          We're committed to building tools that help you ship faster, collaborate
          better, and deliver exceptional products that users love.
        </p>
      </div>
      <div class="bg-slate-100 rounded-3xl p-8 md:p-10">
        <div class="grid grid-cols-2 gap-4">
          <div class="bg-white p-6 rounded-2xl shadow-sm hover:shadow-md transition-shadow duration-300">
            <div class="text-4xl font-bold bg-gradient-to-r from-blue-600 to-blue-500 bg-clip-text text-transparent mb-2">2024</div>
            <div class="text-slate-500 font-medium">Founded</div>
          </div>
          <div class="bg-white p-6 rounded-2xl shadow-sm hover:shadow-md transition-shadow duration-300">
            <div class="text-4xl font-bold bg-gradient-to-r from-green-600 to-green-500 bg-clip-text text-transparent mb-2">150+</div>
            <div class="text-slate-500 font-medium">Team Members</div>
          </div>
          <div class="bg-white p-6 rounded-2xl shadow-sm hover:shadow-md transition-shadow duration-300">
            <div class="text-4xl font-bold bg-gradient-to-r from-purple-600 to-purple-500 bg-clip-text text-transparent mb-2">500+</div>
            <div class="text-slate-500 font-medium">Enterprise Clients</div>
          </div>
          <div class="bg-white p-6 rounded-2xl shadow-sm hover:shadow-md transition-shadow duration-300">
            <div class="text-4xl font-bold bg-gradient-to-r from-orange-600 to-orange-500 bg-clip-text text-transparent mb-2">Global</div>
            <div class="text-slate-500 font-medium">Team</div>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Newsletter Section -->
  <section id="contact" class="py-24 px-6">
    <div class="max-w-4xl mx-auto text-center">
      <div class="inline-flex items-center gap-2 px-4 py-2 bg-green-50 text-green-700 rounded-full text-sm font-semibold mb-6">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        No spam, ever
      </div>
      <h2 class="text-4xl md:text-5xl font-bold text-slate-900 mb-4 tracking-tight">Stay in the loop</h2>
      <p class="text-xl text-slate-600 mb-10 max-w-xl mx-auto">Get the latest updates, product news, and developer tips delivered straight to your inbox.</p>

      {#if subscribed}
        <div class="p-8 bg-gradient-to-r from-green-50 to-emerald-50 border border-green-200 rounded-2xl inline-block animate-fade-in">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-green-500 rounded-full flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <div class="text-left">
              <div class="text-green-900 font-bold text-lg">You're subscribed!</div>
              <div class="text-green-700">Check your inbox for a confirmation email.</div>
            </div>
          </div>
        </div>
      {:else}
        <form onsubmit={(e) => { e.preventDefault(); handleSubscribe(); }} class="flex flex-wrap gap-4 justify-center max-w-lg mx-auto">
          <div class="relative flex-1 min-w-[250px]">
            <svg class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
            <input
              type="email"
              bind:value={email}
              placeholder="Enter your email address"
              class="w-full pl-12 pr-4 py-4 bg-white border-2 border-slate-200 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10 transition-all duration-200"
            />
          </div>
          <button
            type="submit"
            disabled={loading || !email.trim()}
            class="px-8 py-4 bg-blue-600 text-white font-semibold rounded-xl hover:bg-blue-700 hover:shadow-lg hover:-translate-y-0.5 disabled:bg-blue-400 disabled:cursor-not-allowed disabled:hover:translate-0 transition-all duration-200 flex items-center gap-2"
          >
            {#if loading}
              <svg class="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            {/if}
            Subscribe
          </button>
        </form>
        <p class="text-slate-400 text-sm mt-4">Unsubscribe anytime. Read our privacy policy.</p>
      {/if}
    </div>
  </section>

  <!-- Footer -->
  <footer class="bg-white border-t border-slate-200 py-12 px-6">
    <div class="max-w-6xl mx-auto">
      <div class="flex flex-wrap justify-between items-center gap-8">
        <div>
          <div class="text-2xl font-bold bg-gradient-to-r from-slate-900 to-slate-700 bg-clip-text text-transparent mb-2">Prima</div>
          <p class="text-slate-500 text-sm">Building tools for modern teams.</p>
        </div>
        <div class="flex flex-wrap gap-6 md:gap-8 text-slate-600">
          <a href="#features" class="hover:text-slate-900 transition-colors duration-200">Features</a>
          <a href="#about" class="hover:text-slate-900 transition-colors duration-200">About</a>
          <a href="#contact" class="hover:text-slate-900 transition-colors duration-200">Contact</a>
          <a href="#" class="hover:text-slate-900 transition-colors duration-200">Privacy</a>
          <a href="#" class="hover:text-slate-900 transition-colors duration-200">Terms</a>
        </div>
        <div class="text-slate-400 text-sm">
          &copy; 2024 Prima. All rights reserved.
        </div>
      </div>
    </div>
  </footer>
</main>

<style>
  @keyframes fade-in {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  @keyframes fade-in-up {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .animate-fade-in {
    animation: fade-in 0.6s ease-out forwards;
  }

  .animate-fade-in-up {
    animation: fade-in-up 0.6s ease-out forwards;
  }

  .delay-100 {
    animation-delay: 0.1s;
  }

  .delay-200 {
    animation-delay: 0.2s;
  }
</style>
