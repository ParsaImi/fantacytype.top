<script lang="ts">
  import '../app.css';
  import { onMount } from 'svelte';

  let username = '';
  let isAuthenticated = false;

  onMount(() => {
    // Check if user has a valid token
    const token = localStorage.getItem('auth_token');
    if (token) {
      try {
        // Basic validation and extract username from token
        const parts = token.split('.');
        if (parts.length === 3) {
          // Decode the payload (middle part of the JWT)
          const payload = JSON.parse(atob(parts[1]));
          if (payload && payload.sub) {
            username = payload.sub;
            isAuthenticated = true;
          }
        }
      } catch (error) {
        console.error('Error decoding token:', error);
      }
    }
  });
</script>

<div class="app">
  <header>
    <nav>
      <div class="logo">
        <a href="/">Tenfinger</a>
      </div>
      <div class="nav-links">
        <a href="/">Home</a>
        <a href="/game">Play</a>
        
        <div class="auth-section">
          {#if isAuthenticated}
            <span class="username">Welcome, {username}</span>
            <button class="logout-btn" on:click={() => {
              localStorage.removeItem('auth_token');
              isAuthenticated = false;
              username = '';
              window.location.href = '/';
            }}>Logout</button>
          {:else}
            <a href="/login" class="login-btn">Login</a>
          {/if}
        </div>
      </div>
    </nav>
  </header>

  <main>
    <slot />
  </main>

  <footer>
    <p>© {new Date().getFullYear()} Tenfinger Typing Competition</p>
  </footer>
</div>

<style>
  /* Dark Mode Styling - Modern UI based on typing game component */
  :global(body) {
    font-family: 'IRANSans', 'Vazir', 'Tahoma', 'Segoe UI', -apple-system, BlinkMacSystemFont, sans-serif;
    background-color: #121212;
    color: #e0e0e0;
    margin: 0;
    padding: 0;
    min-height: 100vh;
    transition: background-color 0.3s ease;
  }
  
  .app {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }
  
  header {
    background-color: #1e1e1e;
    box-shadow: 0 4px 6px rgba(0,0,0,0.2);
  }
  
  nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
    max-width: 1200px;
    margin: 0 auto;
    padding: 1rem 2rem;
  }
  
  .logo a {
    font-size: 1.5rem;
    font-weight: bold;
    color: #bb86fc; /* Purple accent from typing game */
    text-decoration: none;
    transition: color 0.2s ease;
  }
  
  .logo a:hover {
    color: #c9a0ff; /* Lighter purple on hover */
    text-shadow: 0 0 8px rgba(187, 134, 252, 0.4);
  }
  
  .nav-links {
    display: flex;
    gap: 2rem;
    align-items: center;
  }
  
  .nav-links a {
    color: #e0e0e0;
    text-decoration: none;
    font-weight: 500;
    position: relative;
    padding: 0.3rem 0;
    transition: color 0.2s ease;
  }
  
  .nav-links a:hover {
    color: #bb86fc;
  }
  
  .nav-links a::after {
    content: '';
    position: absolute;
    width: 0;
    height: 2px;
    bottom: 0;
    left: 0;
    background-color: #bb86fc;
    transition: width 0.3s ease;
  }
  
  .nav-links a:hover::after {
    width: 100%;
  }
  
  .auth-section {
    display: flex;
    align-items: center;
    margin-left: 1rem;
    gap: 1rem;
  }
  
  .username {
    color: #bb86fc;
    font-weight: 500;
  }
  
  .logout-btn {
    padding: 0.5rem 1rem;
    font-size: 0.875rem;
    color: white;
    background-color: #e24a4a;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }
  
  .logout-btn:hover {
    background-color: #cf3643;
    transform: translateY(-2px);
    box-shadow: 0 6px 10px rgba(0, 0, 0, 0.15);
  }
  
  .login-btn {
    padding: 0.5rem 1rem !important;
    font-size: 0.875rem;
    color: white !important;
    background-color: #bb86fc !important;
    border: none;
    border-radius: 8px !important;
    cursor: pointer;
    transition: all 0.2s ease !important;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }
  
  .login-btn:hover {
    background-color: #c9a0ff !important;
    transform: translateY(-2px);
    box-shadow: 0 6px 10px rgba(0, 0, 0, 0.15);
  }

  /* Remove the ::after effect for login button */
  .login-btn::after {
    display: none;
  }
  
  main {
    flex: 1;
    max-width: 1200px;
    width: 100%;
    margin: 0 auto;
    padding: 2rem;
  }
  
  footer {
    background-color: #1e1e1e;
    text-align: center;
    padding: 1.5rem;
    margin-top: auto;
    color: #757575;
    box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.1);
  }

  /* Responsive adjustments */
  @media (max-width: 768px) {
    nav {
      flex-direction: column;
      padding: 1rem;
    }
    
    .logo {
      margin-bottom: 1rem;
    }
    
    .nav-links {
      flex-wrap: wrap;
      justify-content: center;
      gap: 1rem;
    }
    
    .auth-section {
      margin-top: 1rem;
      margin-left: 0;
    }
  }
  
  @media (max-width: 480px) {
    .nav-links {
      flex-direction: column;
      align-items: center;
    }
    
    main {
      padding: 1rem;
    }
  }
</style>
