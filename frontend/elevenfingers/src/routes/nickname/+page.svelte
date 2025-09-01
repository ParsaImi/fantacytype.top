<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  
  let nickname = '';
  let error = '';
  let destination = '';
  
  onMount(() => {
    // Check if user already has a nickname
    const existingNickname = localStorage.getItem('user_nickname');
    
    if (existingNickname) {
      nickname = existingNickname;
    }
    
    // Get the destination from URL query parameter or default to /rooms
    const url = new URL(window.location.href);
    destination = url.searchParams.get('redirect') || '/rooms';
  });
  
  function handleSubmit() {
    if (!nickname || nickname.trim().length < 3) {
      error = 'Please enter a nickname (at least 3 characters)';
      return;
    }
    
    if (nickname.trim().length > 20) {
      error = 'Nickname must be 20 characters or less';
      return;
    }
    
    // Store nickname in localStorage
    localStorage.setItem('user_nickname', nickname.trim());
    
    // Navigate to destination
    goto(destination);
  }
</script>

<svelte:head>
  <title>Enter Your Nickname - Tenfinger Typing Game</title>
  <meta name="description" content="Set your nickname for the multiplayer typing game" />
</svelte:head>

<div class="container">
  <div class="nickname-form">
    <h1>Welcome to Tenfinger Typing Game</h1>
    <p class="subtitle">Enter a nickname to continue</p>
    
    <form on:submit|preventDefault={handleSubmit}>
      <div class="form-group">
        <label for="nickname">Your Nickname</label>
        <input 
          type="text" 
          id="nickname" 
          bind:value={nickname} 
          placeholder="Enter a nickname (3-20 characters)"
          class:error={error}
          autocomplete="off"
        />
        {#if error}
          <p class="error-message">{error}</p>
        {/if}
      </div>
      
      <button type="submit" class="continue-btn">Continue</button>
    </form>
  </div>
</div>

<style>
  /* Dark Mode Styling - Modern UI for Nickname Page */
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
    max-width: 600px;
    margin: 0 auto;
    padding: 2rem;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .nickname-form {
    width: 100%;
    background-color: #1e1e1e;
    padding: 2.5rem;
    border-radius: 12px;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.3);
    text-align: center;
    border: 1px solid #333;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    animation: float-in 0.8s ease forwards;
  }
  
  @keyframes float-in {
    0% {
      opacity: 0;
      transform: translateY(20px);
    }
    100% {
      opacity: 1;
      transform: translateY(0);
    }
  }
  
  .nickname-form:hover {
    transform: translateY(-5px);
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.4);
    border-color: #bb86fc;
  }
  
  h1 {
    font-size: 2.25rem;
    margin-bottom: 0.75rem;
    color: #bb86fc;
    text-shadow: 0 2px 8px rgba(187, 134, 252, 0.3);
  }
  
  .subtitle {
    font-size: 1.1rem;
    color: #aaa;
    margin-bottom: 2rem;
  }
  
  .form-group {
    text-align: left;
    margin-bottom: 1.5rem;
  }
  
  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: #e0e0e0;
  }
  
  input {
    width: 100%;
    padding: 0.875rem;
    font-size: 1rem;
    background-color: #2d2d2d;
    border: 2px solid #424242;
    color: #e0e0e0;
    border-radius: 8px;
    transition: all 0.2s ease;
    box-sizing: border-box;
  }
  
  input::placeholder {
    color: #888;
  }
  
  input:focus {
    outline: none;
    border-color: #bb86fc;
    box-shadow: 0 0 0 2px rgba(187, 134, 252, 0.2);
  }
  
  input.error {
    border-color: #cf6679;
    box-shadow: 0 0 0 2px rgba(207, 102, 121, 0.2);
  }
  
  .error-message {
    color: #cf6679;
    margin-top: 0.5rem;
    font-size: 0.875rem;
  }
  
  .continue-btn {
    width: 100%;
    padding: 0.875rem;
    font-size: 1.125rem;
    font-weight: 600;
    color: #121212;
    background-color: #bb86fc;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2);
  }
  
  .continue-btn:hover {
    background-color: #c9a0ff;
    transform: translateY(-2px);
    box-shadow: 0 6px 10px rgba(0, 0, 0, 0.3);
  }
  
  /* Responsive adjustments */
  @media (max-width: 768px) {
    .nickname-form {
      padding: 2rem;
    }
    
    h1 {
      font-size: 2rem;
    }
  }
  
  @media (max-width: 480px) {
    .container {
      padding: 1rem;
    }
    
    .nickname-form {
      padding: 1.5rem;
    }
  }
</style>
