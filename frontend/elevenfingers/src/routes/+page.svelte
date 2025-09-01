<script lang="ts">
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  
  let isAuthenticated = false;
  
  onMount(() => {
    // Check if user has a valid token
    const token = localStorage.getItem('auth_token');
    if (token) {
      try {
        // Validate token by checking if it's a valid JWT structure
        const parts = token.split('.');
        if (parts.length === 3) {
          // Basic validation passed, consider them authenticated
          isAuthenticated = true;
        }
      } catch (error) {
        console.error('Invalid token format');
      }
    }
  });
  
  function startGame() {
    goto("/rooms") 
  }
</script>

<svelte:head>
  <title>Tenfinger Typing Game</title>
  <meta name="description" content="Multiplayer typing game" />
</svelte:head>

<div class="container">
  <div class="hero">
    <h1>Tenfinger Typing Game</h1>
    <p class="subtitle">Compete with players from around the world in real-time typing challenges</p>
    
    <div class="cta-buttons">
      <button class="btn-play" on:click={startGame}>Start Game</button>
      {#if !isAuthenticated}
        <a href="/login" class="btn-main">Login / Register</a>
      {/if}
    </div>
  </div>
  
  <div class="features">
    <div class="feature">
      <h3>Real-Time Competition</h3>
      <p>Race against other players to see who can type the fastest with perfect accuracy</p>
    </div>
    
    <div class="feature">
      <h3>Improve Your Skills</h3>
      <p>Track your progress and see your typing speed and accuracy improve over time</p>
    </div>
    
    <div class="feature">
      <h3>Multiple Rooms</h3>
      <p>Join existing game rooms or create your own custom room to challenge friends</p>
    </div>
  </div>
</div>

<style>
  /* Dark Mode Styling - Modern UI for Landing Page */
  :global(body) {
    font-family: 'IRANSans', 'Vazir', 'Tahoma', 'Segoe UI', -apple-system, BlinkMacSystemFont, sans-serif;
    background-color: #121212;
    color: #e0e0e0;
    text-align: center;
    margin: 0;
    padding: 0;
    min-height: 100vh;
    transition: background-color 0.3s ease;
  }
  
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }
  
  .hero {
    text-align: center;
    padding: 4rem 1rem;
    margin-bottom: 4rem;
    animation: gentle-float 3s infinite ease-in-out;
  }
  
  @keyframes gentle-float {
    0%, 100% { transform: translateY(0); }
    50% { transform: translateY(-10px); }
  }
  
  h1 {
    font-size: 3.5rem;
    margin-bottom: 1.5rem;
    color: #bb86fc; /* Purple accent from the game */
    font-weight: 600;
    letter-spacing: -0.5px;
    text-shadow: 0 2px 8px rgba(187, 134, 252, 0.3);
  }
  
  .subtitle {
    font-size: 1.25rem;
    color: #aaa;
    max-width: 600px;
    margin: 0 auto 3rem auto;
    line-height: 1.6;
  }
  
  .cta-buttons {
    display: flex;
    justify-content: center;
    gap: 1.5rem;
    margin-bottom: 2rem;
    flex-wrap: wrap;
  }
  
  .btn-play, .btn-main {
    padding: 0.875rem 2.5rem;
    font-size: 1.125rem;
    font-weight: 600;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2);
    text-decoration: none;
  }
  
  .btn-play {
    background-color: #bb86fc;
    color: #121212;
    border: none;
  }
  
  .btn-play:hover {
    background-color: #c9a0ff;
    transform: translateY(-2px);
    box-shadow: 0 6px 10px rgba(0, 0, 0, 0.3);
  }
  
  .btn-main {
    background-color: #2d2d2d;
    color: #e0e0e0;
    border: 1px solid #424242;
  }
  
  .btn-main:hover {
    background-color: #3d3d3d;
    transform: translateY(-2px);
    box-shadow: 0 6px 10px rgba(0, 0, 0, 0.3);
  }
  
  .features {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 2rem;
    margin-bottom: 4rem;
  }
  
  .feature {
    background-color: #1e1e1e;
    padding: 2.5rem 2rem;
    border-radius: 12px;
    text-align: center;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    border: 1px solid #333;
  }
  
  .feature:hover {
    transform: translateY(-5px);
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.3);
    border-color: #bb86fc;
  }
  
  .feature h3 {
    margin-top: 0;
    margin-bottom: 1.5rem;
    color: #bb86fc;
    font-size: 1.5rem;
    font-weight: 500;
  }
  
  .feature p {
    color: #aaa;
    margin: 0;
    font-size: 1.1rem;
    line-height: 1.6;
  }
  
  /* Responsive adjustments */
  @media (max-width: 768px) {
    h1 {
      font-size: 2.5rem;
    }
    
    .subtitle {
      font-size: 1.1rem;
    }
    
    .feature {
      padding: 2rem 1.5rem;
    }
  }
  
  @media (max-width: 480px) {
    .cta-buttons {
      flex-direction: column;
      align-items: center;
      gap: 1rem;
    }
    
    .btn-play, .btn-main {
      width: 100%;
      max-width: 280px;
    }
  }
</style>
