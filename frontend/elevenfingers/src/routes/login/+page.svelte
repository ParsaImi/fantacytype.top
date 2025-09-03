<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  
  let username = '';
  let password = '';
  let rememberMe = false;
  let isLoading = false;
  let error: string | null = null;
  let showPassword = false;
  
  // Check if user is already logged in
  onMount(() => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      // Redirect to the main page if user is already authenticated
      goto('/');
    }
  });
  
  async function handleSubmit() {
    // Reset error state
    error = null;
    
    // Validate form
    if (!username || !password) {
      error = 'Please enter both username and password';
      return;
    }
    
    // Set loading state
    isLoading = true;
    try {
      const response = await fetch('https://auth.fantacytype.top/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: new URLSearchParams({
            grant_type: 'password',
            username: username,
            password: password
        }) 
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.message || 'Login failed');
      }
      
      // Store auth token in localStorage
      if (data.access_token) {
        localStorage.setItem('auth_token', data.access_token);
      }
      
      // Redirect to rooms page
      goto('/rooms');
    } catch (err) {
      console.error('Login error:', err);
      error = err.message || 'Failed to login. Please try again.';
    } finally {
      isLoading = false;
    }
  }
  
  function togglePasswordVisibility() {
    showPassword = !showPassword;
  }
  
  function goToRegister() {
    goto('/register');
  }
</script>

<svelte:head>
  <title>Login - Tenfinger Typing Game</title>
  <meta name="description" content="Login to your Tenfinger typing game account" />
</svelte:head>

<div class="container">
  <div class="login-container">
    <h2>Login to Your Account</h2>
    
    {#if error}
      <div class="error-message">
        {error}
      </div>
    {/if}
    
    <form on:submit|preventDefault={handleSubmit}>
      <div class="form-group">
        <label for="username">Username or Email</label>
        <input 
          type="text" 
          id="username" 
          bind:value={username} 
          placeholder="Enter your username or email"
          autocomplete="username"
          disabled={isLoading}
        />
      </div>
      
      <div class="form-group">
        <label for="password">Password</label>
        <div class="password-input-container">
          <input 
            type={showPassword ? 'text' : 'password'} 
            id="password" 
            bind:value={password} 
            placeholder="Enter your password"
            autocomplete="current-password"
            disabled={isLoading}
          />
          <button 
            type="button" 
            class="toggle-password" 
            on:click={togglePasswordVisibility}
            disabled={isLoading}
          >
            {showPassword ? 'Hide' : 'Show'}
          </button>
        </div>
      </div>
      
      <div class="form-options">
        <div class="remember-me">
          <input 
            type="checkbox" 
            id="rememberMe" 
            bind:checked={rememberMe}
            disabled={isLoading}
          />
          <label for="rememberMe">Remember me</label>
        </div>
        
        <a href="/forgot-password" class="forgot-password">Forgot password?</a>
      </div>
      
      <button 
        type="submit" 
        class="login-button"
        disabled={isLoading}
      >
        {isLoading ? 'Logging in...' : 'Login'}
      </button>
    </form>
    
    <div class="register-prompt">
      <p>Don't have an account? <a href="/register" on:click|preventDefault={goToRegister}>Register now</a></p>
    </div>
    
    <div class="back-to-home">
      <a href="/">Back to Home</a>
    </div>
  </div>
</div>

<style>
  /* Dark Mode Styling - Matching the Landing Page */
  :global(body) {
    font-family: 'IRANSans', 'Vazir', 'Tahoma', 'Segoe UI', -apple-system, BlinkMacSystemFont, sans-serif;
    background-color: #121212;
    color: #e0e0e0;
    margin: 0;
    padding: 0;
    min-height: 100vh;
    transition: background-color 0.3s ease;
  }
  
  .container {
    max-width: 1000px;
    margin: 0 auto;
    padding: 2rem;
    display: flex;
    justify-content: center;
    min-height: calc(100vh - 4rem);
  }
  
  .login-container {
    width: 100%;
    max-width: 500px;
    background-color: #1e1e1e;
    padding: 2.5rem;
    border-radius: 12px;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
    border: 1px solid #333;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
  }
  
  .login-container:hover {
    transform: translateY(-5px);
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.3);
    border-color: #bb86fc;
  }
  
  h2 {
    font-size: 2rem;
    margin-bottom: 1.5rem;
    text-align: center;
    color: #bb86fc;
    font-weight: 600;
    letter-spacing: -0.5px;
    text-shadow: 0 2px 8px rgba(187, 134, 252, 0.3);
  }
  
  .form-group {
    margin-bottom: 1.5rem;
  }
  
  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: #e0e0e0;
  }
  
  input[type="text"],
  input[type="password"] {
    width: 100%;
    padding: 0.75rem 1rem;
    border: 1px solid #333;
    border-radius: 8px;
    font-size: 1rem;
    transition: border-color 0.2s;
    background-color: #2d2d2d;
    color: #e0e0e0;
  }
  
  input[type="text"]:focus,
  input[type="password"]:focus {
    border-color: #bb86fc;
    outline: none;
    box-shadow: 0 0 0 2px rgba(187, 134, 252, 0.2);
  }
  
  input[type="text"]::placeholder,
  input[type="password"]::placeholder {
    color: #888;
  }
  
  .password-input-container {
    position: relative;
  }
  
  .toggle-password {
    position: absolute;
    right: 10px;
    top: 50%;
    transform: translateY(-50%);
    background: none;
    border: none;
    color: #aaa;
    font-size: 0.85rem;
    cursor: pointer;
    padding: 0.25rem 0.5rem;
  }
  
  .toggle-password:hover {
    color: #bb86fc;
  }
  
  .form-options {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    font-size: 0.9rem;
  }
  
  .remember-me {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  
  .forgot-password {
    color: #bb86fc;
    text-decoration: none;
  }
  
  .forgot-password:hover {
    text-decoration: underline;
  }
  
  .login-button {
    width: 100%;
    padding: 0.875rem;
    background-color: #bb86fc;
    color: #121212;
    border: none;
    border-radius: 8px;
    font-size: 1.125rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2);
  }
  
  .login-button:hover {
    background-color: #c9a0ff;
    transform: translateY(-2px);
    box-shadow: 0 6px 10px rgba(0, 0, 0, 0.3);
  }
  
  .login-button:disabled {
    background-color: #6b4b91;
    cursor: not-allowed;
    transform: none;
    box-shadow: none;
  }
  
  .error-message {
    background-color: rgba(211, 47, 47, 0.1);
    border: 1px solid #d32f2f;
    color: #ff6b6b;
    padding: 0.75rem 1rem;
    border-radius: 8px;
    margin-bottom: 1.5rem;
    font-size: 0.9rem;
  }
  
  .register-prompt {
    text-align: center;
    margin-top: 2rem;
    padding-top: 1.5rem;
    border-top: 1px solid #333;
    color: #aaa;
  }
  
  .register-prompt a {
    color: #bb86fc;
    font-weight: 500;
    text-decoration: none;
  }
  
  .register-prompt a:hover {
    text-decoration: underline;
  }
  
  .back-to-home {
    text-align: center;
    margin-top: 1.5rem;
    font-size: 0.9rem;
  }
  
  .back-to-home a {
    color: #aaa;
    text-decoration: none;
  }
  
  .back-to-home a:hover {
    color: #bb86fc;
    text-decoration: underline;
  }
  
  /* Responsive adjustments */
  @media (max-width: 768px) {
    .login-container {
      padding: 2rem 1.5rem;
    }
    
    h2 {
      font-size: 1.8rem;
    }
  }
  
  @media (max-width: 480px) {
    .form-options {
      flex-direction: column;
      align-items: flex-start;
      gap: 1rem;
    }
    
    .forgot-password {
      align-self: flex-end;
    }
  }
</style>
